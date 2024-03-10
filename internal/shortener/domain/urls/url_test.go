package urls_test

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	url "github.com/Leopold1975/url_shortener/internal/shortener/domain/urls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tc = []struct {
	in  string
	exp bool
}{
	{
		in:  "https://youtube.com",
		exp: true,
	},
	{
		in:  "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		exp: true,
	},
	{
		in:  "http://www.youtube.com/watch?v=dQw4w9WgXcQ",
		exp: true,
	},
	{
		in:  "www.youtube.com/watch?v=dQw4w9WgXcQ",
		exp: false,
	},
	{
		in:  "http://www.youtube.com/watch?v=dQw4w9WgXcQ\t",
		exp: false,
	},
}

func TestValidate(t *testing.T) {
	for _, c := range tc {
		assert.Equal(t, c.exp, url.Validate(c.in))
	}
}

var tests = []string{
	"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	"https://github.com/openai/gpt-3",
	"https://www.youtube.com/watch?v=3tmd-ClpJxA",
	"https://github.com/microsoft/TypeScript",
	"https://www.youtube.com/watch?v=jNQXAC9IVRw",
	"https://github.com/public-apis/public-apis",
	"https://www.youtube.com/watch?v=ktvTqknDobU",
	"https://github.com/flutter/flutter",
	"https://www.youtube.com/watch?v=7xZY8VJHqU4",
	"https://github.com/donnemartin/system-design-primer",
	"https://www.youtube.com/watch?v=2vjPBrBU-TM",
	"https://github.com/github/gitignore",
	"https://github.com/ohmyzsh/ohmyzsh",
	"https://www.youtube.com/watch?v=Bmz67ErIRa4",
	"https://github.com/vinta/awesome-python",
	"https://www.youtube.com/watch?v=9bZkp7q19f0",
	"https://github.com/airbnb/javascript",
	"https://www.youtube.com/watch?v=YQHsXMglC9A",
	"https://github.com/apache/kafka",
	"https://www.youtube.com/watch?v=ktvTqknDobU",
	"https://github.com/puppeteer/puppeteer",
	"https://www.youtube.com/watch?v=3tmd-ClpJxA",
	"https://www.youtube.com/watch?v=2vjPBrBU-TM",
	"https://github.com/axios/axios",
	"https://www.youtube.com/watch?v=ktvTqknDobU",
	"https://github.com/tensorflow/models",
	"https://www.youtube.com/watch?v=fRh_vgS2dFE",
	"https://github.com/toddmotto/public-apis",
	"https://www.youtube.com/watch?v=ktvTqknDobU",
	"https://github.com/FortAwesome/Font-Awesome",
	"https://www.youtube.com/watch?v=BGBM5vWiBLo",
	"https://github.com/moby/moby",
	"https://www.youtube.com/watch?v=7zp1TbLFPp8",
	"https://github.com/tensorflow/tensorflow",
	"https://github.com/apple/swift",
	"https://www.youtube.com/watch?v=ktvTqknDobU",
	"https://github.com/microsoft/terminal",
	"https://www.youtube.com/watch?v=3tmd-ClpJxA",
	"https://github.com/facebook/react",
	"https://www.youtube.com/watch?v=9bZkp7q19f0",
	"https://github.com/getify/You-Dont-Know-JS",
	"https://github.com/facebook/react-native",
	"https://www.youtube.com/watch?v=YQHsXMglC9A",
	"https://github.com/jwasham/coding-interview-university",
	"https://www.youtube.com/watch?v=F-1weFCiYBA",
	"https://github.com/keras-team/keras",
	"https://github.com/EbookFoundation/free-programming-books",
}

func TestBasic(t *testing.T) {
	m := make(map[string]int, len(tests))
	m1 := make(map[string]int, len(tests))
	for _, tc := range tests {
		m1[tc]++
		if m1[tc] > 1 {
			continue
		}
		u, err := url.PrepareURL(tc)
		require.NoError(t, err)

		assert.Equal(t, tc, u.LongURL)
		m[u.ShortURL]++
		assert.Equal(t, 1, m[u.ShortURL])
	}
}

func TestDurability(t *testing.T) {
	content, err := os.ReadFile("testcases")
	require.NoError(t, err)
	m := make(map[string]int, len(tests))
	m1 := make(map[string]int, len(tests))

	s := bufio.NewScanner(bytes.NewReader(content))
	var u string
	for s.Scan() {
		if s.Err() != nil {
			break
		}
		u = s.Text()
		m1[u]++
		if m1[u] > 1 {
			continue
		}

		ur, err := url.PrepareURL(u)
		require.NoError(t, err)

		assert.Equal(t, u, ur.LongURL)
		m[ur.ShortURL]++
		assert.Equal(t, 1, m[ur.ShortURL])
	}
}
