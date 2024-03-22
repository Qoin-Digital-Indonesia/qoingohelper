package keycloak

import (
	"context"
	"fmt"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

type Client struct {
	Client           *gocloak.GoCloak
	ClientProperties *ClientProperties
}

type ClientProperties struct {
	Realm        string
	ClientId     string
	ClientSecret string
}

func NewClient(urlKeycloak, realm, clientId, clientSecret string) *Client {
	client := gocloak.NewClient(urlKeycloak)
	return &Client{
		Client: client,
		ClientProperties: &ClientProperties{
			Realm:        realm,
			ClientId:     clientId,
			ClientSecret: clientSecret,
		},
	}
}

func Login() string {
	client := gocloak.NewClient(os.Getenv("BaseUrl"))
	ctx := context.Background()
	token, err := client.Login(ctx, os.Getenv("ClientId"), os.Getenv("ClientSecret"), os.Getenv("Realm"), "ariegeorgee", "ari")
	if err != nil {
		fmt.Println(err.Error())
		panic("Something wrong with the credentials or url")
	}

	fmt.Println("TOKEN : ", token.AccessToken)
	fmt.Println("EXPIRE IN : ", token.ExpiresIn)
	// fmt.Println("IDTOKEN : ", token.IDToken)
	// fmt.Println("BEFORE POLICY : ", token.NotBeforePolicy)
	// fmt.Println("REFRESH TOKEN EXPIRE : ", token.RefreshExpiresIn)
	fmt.Println("REFRESH TOKEN : ", token.RefreshToken)
	// fmt.Println("SCOPE : ", token.Scope)
	// fmt.Println("SESSION STATE : ", token.SessionState)
	// fmt.Println("TOKEN TYPE : ", token.TokenType)

	return token.RefreshToken
}

func RefreshToken(refreshToken string) string {
	client := gocloak.NewClient(os.Getenv("BaseUrl"))
	ctx := context.Background()
	token, err := client.RefreshToken(ctx, refreshToken, os.Getenv("ClientId"), os.Getenv("ClientSecret"), os.Getenv("Realm"))
	if err != nil {
		fmt.Println(err.Error())
		panic("Something wrong with the credentials or url")
	}

	fmt.Println("================ REFRESH TOKEN VALID ================")

	fmt.Println("TOKEN : ", token.AccessToken)
	fmt.Println("EXPIRE IN : ", token.ExpiresIn)
	fmt.Println("IDTOKEN : ", token.IDToken)
	fmt.Println("BEFORE POLICY : ", token.NotBeforePolicy)
	fmt.Println("REFRESH TOKEN EXPIRE : ", token.RefreshExpiresIn)
	fmt.Println("REFRESH TOKEN : ", token.RefreshToken)
	fmt.Println("SCOPE : ", token.Scope)
	fmt.Println("SESSION STATE : ", token.SessionState)
	fmt.Println("TOKEN TYPE : ", token.TokenType)

	return token.AccessToken
}

func VerifyToken(token string) bool {
	// client := gocloak.NewClient(os.Getenv("BaseUrl"))
	// ctx := context.Background()
	// jToken, jClaim, err := client.DecodeAccessToken(ctx, token, os.Getenv("Realm"))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic("Something wrong with the credentials or url")

	// }

	// fmt.Println("CLAIMS : ", jToken.Claims)
	// fmt.Println("CLAIMS : ", jClaim)
	// fmt.Println("IS VALID : ", jToken.Valid)

	// jsonString, _ := json.Marshal(jClaim)
	// fmt.Println("JSON CLAIMS : ", string(jsonString))

	// validToken, err := client.RetrospectToken(ctx, token, os.Getenv("ClientId"), os.Getenv("ClientSecret"), os.Getenv("Realm"))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic("Something wrong with the credentials or url")

	// }

	// fmt.Println("VAlID : ", validToken)
	// fmt.Println("VAlID TOKEN : ", validToken.String())
	// fmt.Println("PERMISSIONS : ", validToken.Permissions)
	return true
}

func GetPermissions(token string) {
	client := gocloak.NewClient(os.Getenv("BaseUrl"))
	ctx := context.Background()
	Permissions, err := client.GetUserPermissions(ctx, token, os.Getenv("Realm"), gocloak.GetUserPermissionParams{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Something wrong with the credentials or url")

	}
	fmt.Println("PERMISSIONS : ", Permissions)
}

func GetGroups(token string) {
	client := gocloak.NewClient(os.Getenv("BaseUrl"))
	ctx := context.Background()
	group, err := client.GetGroups(ctx, token, os.Getenv("Realm"), gocloak.GetGroupsParams{})

	client.group
	if err != nil {
		fmt.Println(err.Error())
		panic("Something wrong with the credentials or url")

	}
	fmt.Println("GROUP : ", group)
}
