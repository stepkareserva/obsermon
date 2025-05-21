package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"

	"github.com/stepkareserva/obsermon/internal/server/http/errors"
	"go.uber.org/zap"
)

// header HashSHA256 is forbidden by checker
var signHeader = http.CanonicalHeaderKey("HashSHA256")

// create middleware for check request signature
func Sign(secretkey string, log *zap.Logger) Middleware {
	ev := errors.NewErrorsWriter(log)
	return func(next http.Handler) http.Handler {
		singing := func(w http.ResponseWriter, r *http.Request) {
			// skip non signed messages
			if _, ok := r.Header[signHeader]; !ok {
				next.ServeHTTP(w, r)
				return
			}

			// check request sign
			signOK, err := checkSign(r, secretkey)
			if err != nil {
				ev.WriteError(w, errors.ErrInternalServerError, err.Error())
				return
			}
			if !signOK {
				ev.WriteError(w, errors.ErrInvalidRequestSign)
				return
			}

			// set responce sing
			bw := withSigning(w, secretkey, log)
			next.ServeHTTP(bw, r)
			bw.FlushToClient()

		}
		return http.HandlerFunc(singing)
	}
}

func checkSign(r *http.Request, secretkey string) (bool, error) {
	headerSign, err := hex.DecodeString(r.Header.Get(signHeader))
	if err != nil {
		return false, fmt.Errorf("decoding sign: %w", err)
	}

	// read body and write it back
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return false, fmt.Errorf("reading request body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewReader(body))

	hash := hmac.New(sha256.New, []byte(secretkey))
	if _, err := hash.Write(body); err != nil {
		return false, fmt.Errorf("hash write: %w", err)
	}
	sign := hash.Sum(nil)

	return hmac.Equal(sign, headerSign), nil
}

func withSigning(w http.ResponseWriter, secretkey string, log *zap.Logger) *signingWriter {
	if log == nil {
		log = zap.NewNop()
	}

	return &signingWriter{
		ResponseWriter: w,
		hash:           hmac.New(sha256.New, []byte(secretkey)),
		log:            log,
	}
}

type signingWriter struct {
	http.ResponseWriter
	hash hash.Hash
	log  *zap.Logger
}

var _ http.ResponseWriter = (*signingWriter)(nil)

func (w *signingWriter) Write(data []byte) (int, error) {
	if _, err := w.hash.Write(data); err != nil {
		return 0, fmt.Errorf("data hash: %w", err)
	}
	return w.ResponseWriter.Write(data)
}

func (w *signingWriter) FlushToClient() {
	hash := w.hash.Sum(nil)
	hashString := hex.EncodeToString(hash)
	w.ResponseWriter.Header().Set(signHeader, hashString)
	w.hash.Reset()
}
