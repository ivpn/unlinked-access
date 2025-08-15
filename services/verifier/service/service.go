package service

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	ksmconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/jasonlvhit/gocron"
	"ivpn.net/auth/services/verifier/client/http"
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

type Store interface {
	GetSubscriptions() ([]model.Subscription, error)
	UpdateSubscription(model.Subscription) error
	UpdateSubscriptions([]model.Subscription) error
}

type Service struct {
	Cfg       config.Config
	Store     Store
	Http      http.Http
	KsmClient *kms.Client
}

func New(cfg config.Config, store Store) (*Service, error) {
	ctx := context.Background()
	ksmCfg, err := ksmconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Service{
		Cfg:   cfg,
		Store: store,
		Http: http.Http{
			Cfg: cfg.API,
		},
		KsmClient: kms.NewFromConfig(ksmCfg),
	}, nil
}

func (s *Service) Start() error {
	log.Println("verifier service started")

	err := gocron.Every(1).Hour().Do(s.SyncManifest)
	if err != nil {
		log.Printf("error syncing manifest: %v", err)
	}

	// Start all the pending jobs
	<-gocron.Start()

	return err
}

func (s *Service) SyncManifest() error {
	log.Println("syncing manifest...")
	m, err := s.GetManifest()
	if err != nil {
		return err
	}

	err = s.VerifyManifest(m)
	if err != nil {
		return err
	}

	err = s.UpdateSubscriptions(m)
	if err != nil {
		return err
	}

	log.Printf("manifest synced successfully: %v", m.ID)

	return nil
}

func (s *Service) GetManifest() (model.Manifest, error) {
	manifest, err := s.Http.GetManifest()
	if err != nil {
		log.Printf("error fetching manifest: %v", err)
		return model.Manifest{}, err
	}

	return manifest, nil
}

func (s *Service) VerifyManifest(m model.Manifest) error {
	// TODO: Implement HSM verification
	log.Printf("verifying manifest: %v", m.ID)

	if m.ValidUntil.Before(time.Now()) {
		log.Printf("manifest is expired: %v", m.ValidUntil)
		return fmt.Errorf("manifest is expired")
	}

	signature := m.Signature
	m.Signature = ""

	data, err := json.Marshal(m)
	if err != nil {
		log.Println("error marshaling manifest for signing:", err)
		return err
	}

	digest := sha256.Sum256(data)
	digestBase64 := base64.StdEncoding.EncodeToString(digest[:])

	if s.Cfg.Service.Mock {
		hash512 := sha512.Sum512([]byte(digestBase64))
		digestBase64 = base64.StdEncoding.EncodeToString(hash512[:])

		if digestBase64 != signature {
			log.Printf("manifest signature (mock) does not match: %v != %v", digestBase64, signature)
			return fmt.Errorf("invalid manifest signature (mock)")
		}

		log.Println("manifest signature (mock) OK")

		return nil
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		log.Printf("error decoding signature: %v", err)
		return err
	}

	verifyInput := &kms.VerifyInput{
		KeyId:            &s.Cfg.Service.KeyId,
		Message:          digest[:],
		MessageType:      types.MessageTypeDigest,
		Signature:        sigBytes,
		SigningAlgorithm: types.SigningAlgorithmSpecRsassaPssSha256,
	}

	verifyOut, err := s.KsmClient.Verify(context.Background(), verifyInput)
	if err != nil {
		log.Printf("error verifying manifest signature: %v", err)
		return err
	}
	if !verifyOut.SignatureValid {
		log.Printf("manifest signature is invalid")
		return fmt.Errorf("manifest signature is invalid")
	}

	log.Println("manifest signature OK")

	return nil
}

func (s *Service) UpdateSubscriptions(m model.Manifest) error {
	subs, err := s.Store.GetSubscriptions()
	if err != nil {
		log.Printf("error fetching subscriptions: %v", err)
		return err
	}

	for i, sub := range subs {
		updatedSub, err := UpdateSubscriptionFromManifest(sub, m.Subscriptions)
		if err != nil {
			log.Printf("error updating subscription: %v", err)
			continue
		}

		subs[i] = updatedSub
	}

	err = s.Store.UpdateSubscriptions(subs)
	if err != nil {
		log.Printf("error saving updated subscriptions: %v", err)
		return err
	}

	return nil
}

func UpdateSubscriptionFromManifest(sub model.Subscription, manifestSubs []model.Subscription) (model.Subscription, error) {
	for _, s := range manifestSubs {
		if sub.TokenHash == s.TokenHash {
			sub.IsActive = s.IsActive
			sub.ActiveUntil = s.ActiveUntil
			sub.Tier = s.Tier
			return sub, nil
		}
	}

	return model.Subscription{}, fmt.Errorf("subscription with TokenHash %s not found", sub.TokenHash)
}
