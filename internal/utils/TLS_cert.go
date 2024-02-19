package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"

	"github.com/0loff/go_link_shortener/internal/logger"
	"go.uber.org/zap"
)

// TLSCertCreate - This method for creating a TLS certificate
func TLSCertCreate(certFile, keyFile string) error {
	// Создаем шаблон сертификата
	cert := &x509.Certificate{
		// Указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// Заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// Разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// Сертификат верен начиная с времени создания
		NotBefore: time.Now(),
		// Время жизни сертификата - 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// Устанавливаем использование ключа для цифровой подписи,
		// а так же клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// Создаем новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.Log.Error("Unable to create private key for TLS cert", zap.Error(err))
		return err
	}

	// Создаем сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		logger.Log.Error("Unable to create TLS cert", zap.Error(err))
		return err
	}

	err = PEMFileCreate(certFile, "CERTIFICATE", certBytes)
	if err != nil {
		return err
	}

	err = PEMFileCreate(keyFile, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))
	if err != nil {
		return err
	}

	// Кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	// var certPEM bytes.Buffer
	// pem.Encode(&certPEM, &pem.Block{
	// 	Type:  "CERTIFICATE",
	// 	Bytes: certBytes,
	// })

	// var privateKeyPEM bytes.Buffer
	// pem.Encode(&privateKeyPEM, &pem.Block{
	// 	Type:  "RSA PRIVATE KEY",
	// 	Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	// })

	return nil
}
