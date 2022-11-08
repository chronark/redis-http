package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"reflect"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/recover"
)

type server struct {
	app *fiber.App
	rdb *RDB
}

func newServer(rdb *RDB) *server {

	appConfig := fiber.Config{
		DisableStartupMessage: true,
	}

	app := fiber.New(appConfig)
	// app.Use(recover.New())
	app.Server().LogAllErrors = true
	app.Use(func(c *fiber.Ctx) error {

		err := c.Next()

		if err != nil && c.Response().StatusCode() >= 300 {
			log.Println("[ERROR]", err)
		} else {
			log.Println("[INFO]", c.Method(), c.Response().StatusCode())
		}

		return err
	})

	s := &server{app, rdb}
	s.init()

	return s
}
func (s *server) Listen(addr string) error {
	return s.app.Listen(addr)
}

type resultResponse struct {
	Result any `json:"result"`
}
type errorResponse struct {
	Err string `json:"error,omitempty"`
}

func (r *resultResponse) encode() {
	if r.Result == nil {
		return
	}
	enc := base64.StdEncoding
	t := reflect.TypeOf(r.Result)

	switch t.Kind() {
	case reflect.String:
		r.Result = enc.EncodeToString([]byte(r.Result.(string)))
	case reflect.Slice:
		for i, v := range r.Result.([]any) {
			if v == nil {
				continue
			}
			if reflect.TypeOf(v).Kind() == reflect.String {
				r.Result.([]any)[i] = enc.EncodeToString([]byte(v.(string)))
			}
		}

	}

}

func (s *server) init() {
	s.app.Post("/", func(c *fiber.Ctx) error {
		req := []any{}
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(400).JSON(errorResponse{err.Error()})
		}
		log.Println("Request:", req)
		result, err := s.rdb.call(c.Context(), req...)
		res := &resultResponse{}

		if err != nil && err != redis.Nil {
			return c.Status(500).JSON(errorResponse{fmt.Sprintf("redis error: %s", err.Error())})
		}

		res.Result = result

		if c.Get("Upstash-Encoding") == "base64" {
			res.encode()
		}

		return c.JSON(res)

	})

	s.app.Post("/pipeline", func(c *fiber.Ctx) error {
		req := [][]any{}
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(400).JSON(errorResponse{err.Error()})
		}
		p := s.rdb.client.Pipeline()
		cmds := []*redis.Cmd{}
		for _, args := range req {
			cmds = append(cmds, p.Do(c.Context(), args...))

		}
		_, err = p.Exec(c.Context())
		if err != nil {
			return c.Status(500).JSON(errorResponse{err.Error()})
		}
		result := make([]any, len(cmds))
		for i, cmd := range cmds {
			result[i] = cmd.Val()
		}

		res := &resultResponse{result}
		if c.Get("Upstash-Encoding") == "base64" {
			res.encode()
		}

		return c.JSON(res)

	})

	s.app.Post("/multi-exec", func(c *fiber.Ctx) error {
		req := [][]any{}
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(400).JSON(errorResponse{err.Error()})
		}
		p := s.rdb.client.TxPipeline()
		cmds := []*redis.Cmd{}
		for _, args := range req {
			cmds = append(cmds, p.Do(c.Context(), args...))

		}
		_, err = p.Exec(c.Context())
		if err != nil {
			return c.Status(500).JSON(&errorResponse{err.Error()})
		}
		result := make([]any, len(cmds))
		for i, cmd := range cmds {
			result[i] = cmd.Val()
		}

		res := &resultResponse{result}
		if c.Get("Upstash-Encoding") == "base64" {
			res.encode()
		}

		return c.JSON(res)

	})

}
