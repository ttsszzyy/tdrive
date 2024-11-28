package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

var (
	// NoRedisPass it's default pass value, if you don't want to use pass, you can use it.
	NoRedisPass = ""

	// NodeMode if you use redis, you need to set the variable. This is single node mode.
	NodeMode = redis.NodeType
	// ClusterMode if you use redis, you can use the variable to set redis connection mode. This is cluster mode, more node.
	ClusterMode = redis.ClusterType

	cacheKeyPrefix = "jwt-token:"
)

var (
	ErrNoCache         = errors.New("caching is not enabled")
	ErrInvalidToken    = errors.New("invalid token")
	ErrDisableBacklist = errors.New("disable backlist")
	ErrTokenBacklist   = errors.New("token is in backlist")
)

type (
	JWT interface {
		Token(payload Payload) (string, error)
		Parse(tk string) (Payload, error)
		Store(ctx context.Context, key string, tk string) error
		Load(ctx context.Context, key string) (string, error)
		Discard(tk string) error
		DiscardWithDelete(tk string) error
		DeleteByUid(uid string) error
	}

	jwtEx struct {
		secret    []byte
		expire    int64
		blacklist bool // backlist , default is false
		cache     *redis.Redis
	}

	JwtOption func(*jwtEx)

	// Payload User-defined data
	Payload map[string]interface{}
)

// NewJWT create jwt service
// secret : The secret key
// expire : Validity period of the token
// opts : This is a option param
func NewJWT(secret string, expire int64, opts ...JwtOption) JWT {
	jwtObj := &jwtEx{
		secret: []byte(secret),
		expire: expire,
	}

	for _, opt := range opts {
		opt(jwtObj)
	}

	return jwtObj
}

// SetRedisOpt set redis option
// host : redis host
// mode : connect redis mode, value is node|cluster, defual is node
//
//	node is single, cluster is mutile node
//
// pwd : pwd is password, it required to connect redis, if no pass, you can use variable jwt.NoRedisPass
// func SetRedisOpt(host, mode, pwd string) JwtOption {
// 	if mode != NodeMode || mode != ClusterMode {
// 		mode = NodeMode
// 	}
// 	return func(j *jwtEx) {
// 		if mode == ClusterMode {
// 			j.cache = redis.New(host, func(r *redis.Redis) {
// 				r.Type = mode
// 				r.Pass = pwd
// 			}, redis.Cluster())
// 		} else {
// 			j.cache = redis.New(host, func(r *redis.Redis) {
// 				r.Type = mode
// 				r.Pass = pwd
// 			})
// 		}

// 	}
// }

// SetRedis receive a created redis node client or cluster client directly
func SetRedis(r *redis.Redis) JwtOption {
	return func(j *jwtEx) {
		if j.cache != nil {
			return
		}
		j.cache = r
	}
}

// SetBlackListOpt set blacklist status
// status : true - Enables the blacklist function
//
//	false - Disables the blacklist function
func SetBlackListOpt(status bool) JwtOption {
	return func(je *jwtEx) {
		je.blacklist = status
	}
}

// Token create jwt-token
// payload : this is your custom data
func (j *jwtEx) Token(payload Payload) (string, error) {
	tNow := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = tNow + j.expire
	claims["iat"] = tNow
	for k, v := range payload {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString(j.secret)
}

// GetPayload Parses the token to obtain payload data
// tk : Token to be parsed.
func (j *jwtEx) Parse(tk string) (Payload, error) {

	var claims jwt.MapClaims
	t, err := jwt.ParseWithClaims(tk, &claims, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, errors.New("valid token")
	}

	payload := make(Payload)
	for k, v := range claims {
		payload[k] = v
	}

	return payload, nil
}

// Discard discard jwt token
// tk : your token
// Remark : the function can only be used if the cache is enabled
func (j *jwtEx) Discard(tk string) error {
	if !j.blacklist {
		return ErrDisableBacklist
	}

	if j.cache == nil {
		return ErrNoCache
	}

	data, err := j.Parse(tk)
	if err != nil {
		return nil
	}

	if data == nil {
		return nil
	}

	if _, ok := data["exp"]; !ok {
		return nil
	}
	exp := int64(data["`exp`"].(float64))
	tNow := time.Now().Unix()
	if exp >= tNow {
		return nil
	}

	return j.cache.Setex(tk, "1", int(tNow-exp)+10)
}

// DeleteByUid discard jwt token by user
// uid: default token carrier
func (j *jwtEx) DeleteByUid(uid string) error {
	if !j.blacklist {
		return ErrDisableBacklist
	}

	if j.cache == nil {
		return ErrNoCache
	}
	_, err := j.cache.DelCtx(context.Background(), cacheKey(uid))
	return err
}

// Discard token directly delete it, key is parse of uid
func (j *jwtEx) DiscardWithDelete(tk string) error {

	if !j.blacklist {
		return ErrDisableBacklist
	}

	if j.cache == nil {
		return ErrNoCache
	}
	data, err := j.Parse(tk)

	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}

	uid, ok := data["uid"].(string)
	if !ok {
		return nil
	}

	_, err = j.cache.DelCtx(context.Background(), cacheKey(uid))
	if err != nil {
		return err
	}
	return nil
}

// Store write into storage, need to config redis
func (j *jwtEx) Store(ctx context.Context, key string, tk string) error {
	if j.cache == nil {
		return ErrNoCache
	}

	return j.cache.SetexCtx(ctx, cacheKey(key), tk, int(j.expire))
}

// Store load from storage, need to config redis
func (j *jwtEx) Load(ctx context.Context, key string) (string, error) {
	if j.cache == nil {
		return "", ErrNoCache
	}

	return j.cache.GetCtx(ctx, cacheKey(key))
}

func (j *jwtEx) checkBacklist(tk string) (bool, error) {
	if j.cache == nil {
		return false, ErrNoCache
	}

	val, err := j.cache.Get(tk)
	if err != nil {
		return false, err
	}

	return val == "1", nil
}

func cacheKey(val interface{}) string {
	return fmt.Sprintf("%s%v", cacheKeyPrefix, val)
}
