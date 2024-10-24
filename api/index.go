package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/trashhalo/reddit-rss/pkg/client"
	"golang.org/x/oauth2"
)

// Handler 是Vercel Serverless Function的入口点
func Handler(w http.ResponseWriter, r *http.Request) {
	// 创建Reddit客户端
	httpClient := http.DefaultClient
	var token *oauth2.Token
	baseApiURL := "https://www.reddit.com"

	// 配置OAuth (如果提供了环境变量)
	oauthClientID := os.Getenv("OAUTH_CLIENT_ID")
	if oauthClientID != "" {
		oauthClientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
		oauthCfg := &oauth2.Config{
			ClientID:     oauthClientID,
			ClientSecret: oauthClientSecret,
			Endpoint: oauth2.Endpoint{
				TokenURL:  "https://www.reddit.com/api/v1/access_token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		}

		// 使用Reddit用户名密码登录
		username := os.Getenv("REDDIT_USERNAME")
		password := os.Getenv("REDDIT_PASSWORD")
		token, _ = oauthCfg.PasswordCredentialsToken(r.Context(), username, password)
		if token != nil {
			baseApiURL = "https://oauth.reddit.com"
		}
	}

	// 设置User Agent
	userAgent := os.Getenv("USER_AGENT")
	if userAgent == "" {
		userAgent = "reddit-rss 1.0"
	}

	// 创建Reddit客户端
	redditClient := &client.RedditClient{
		HttpClient: httpClient,
		Token:      token,
		UserAgent:  userAgent,
	}

	// 使用项目的RSS处理器处理请求
	client.RssHandler(baseApiURL, time.Now, redditClient, client.GetArticle, w, r)
}
