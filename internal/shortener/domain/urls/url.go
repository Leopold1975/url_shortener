package urls

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type URL struct {
	UUID      string    `json:"uuid"`
	LongURL   string    `json:"longUrl"`
	ShortURL  string    `json:"shortUrl"`
	CreatedAt time.Time `json:"createdAt"`
	Clicks    int64     `json:"clicks"`
}

func PrepareURL(longURL string) (URL, error) {
	uid := uuid.New()

	short, err := getShort(longURL)
	if err != nil {
		return URL{}, err
	}

	return URL{
		UUID:      uid.String(),
		LongURL:   longURL,
		ShortURL:  short,
		CreatedAt: time.Now(),
	}, nil
}

func Validate(longURL string) bool {
	_, err := url.ParseRequestURI(longURL)

	return err == nil
}

func getShort(longURL string) (string, error) {
	// This was really inefficient way. Short URL took more symbols, than original one
	// zw := zlib.NewWriter(&buf)
	// short := base64.RawURLEncoding.EncodeToString(buf.Bytes())

	h := md5.Sum([]byte(longURL + strconv.Itoa(int(time.Now().Unix()))))

	short := hex.EncodeToString(h[:6])

	return short, nil
}

// The Previous version.
// func GetLongUrl(shortURL string) (string, error) {
// 	dec, err := base64.RawURLEncoding.DecodeString(shortURL)
// 	if err != nil {
// 		return "", err
// 	}

// 	var buf bytes.Buffer
// 	if _, err = buf.Write(dec); err != nil {
// 		return "", err
// 	}

// 	r, err := zlib.NewReader(&buf)
// 	if err != nil {
// 		return "", err
// 	}
// 	r.Close()

// 	buf = bytes.Buffer{}
// 	io.Copy(&buf, r)
// 	return buf.String(), nil
// }
