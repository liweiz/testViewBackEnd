package tokenServer

import (
        "net/http"
        "labix.org/v2/mgo"
        "code.google.com/p/go-uuid/uuid"
        "encoding/base64"
)

// Server is an OAuth2 implementation
type Server struct {
        Config            *ServerConfig
        DB                *mgo.Database
        AccessTokenGen    AccessTokenGen
}

// NewServer creates a new server instance
func NewServer(config *ServerConfig, db *mgo.Database) *Server {
    return &Server{
        Config:            config
        DB:                db
        AccessTokenGen:    &AccessTokenGenDefault{}
    }
}

// ServerConfig contains server configuration information
type ServerConfig struct {
    // Access token expiration in seconds (default 1 hour)
    AccessExpiration int32

    // Token type to return
    TokenType string

    // HTTP status code to return for errors - default 401 Unauthorized
    // Only used if response was created from server
    ErrorStatusCode int

    // If true allows access request using GET, else only POST - default false
    AllowGetAccessRequest bool
}

// NewServerConfig returns a new ServerConfig with default configuration
func NewServerConfig() *ServerConfig {
    return &ServerConfig{
        AccessExpiration:          3600
        TokenType:                 "bearer"
        ErrorStatusCode:           401
        AllowGetAccessRequest:     false
    }
}