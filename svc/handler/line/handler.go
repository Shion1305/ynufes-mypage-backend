package line

import (
	"github.com/gin-gonic/gin"
	linePkg "ynufes-mypage-backend/pkg/line"
	"ynufes-mypage-backend/pkg/setting"
	"ynufes-mypage-backend/svc/pkg/domain/query"
	lineDomain "ynufes-mypage-backend/svc/pkg/domain/service/line"
)

type LineAuth struct {
	verifier lineDomain.AuthVerifier
	query    query.User
}

func NewLineAuth() LineAuth {
	config := setting.Get()
	return LineAuth{
		verifier: linePkg.NewAuthVerifier(
			config.ThirdParty.LineLogin.CallbackURI,
			config.ThirdParty.LineLogin.ClientID,
			config.ThirdParty.LineLogin.ClientSecret),
	}
}

func (a LineAuth) VerificationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Request.URL.Query().Get("code")
		state := c.Request.URL.Query().Get("state")
		_, err := a.verifier.RequestAccessToken(code, state)
		if err != nil {
			return
		}
	}
}

func (a LineAuth) StateIssuer() gin.HandlerFunc {
	return func(c *gin.Context) {
		state := a.verifier.IssueNewState()
		_, _ = c.Writer.WriteString(state)
		c.Status(200)
	}
}
