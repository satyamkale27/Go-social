package db

import (
	"context"
	"fmt"
	store2 "github.com/satyamkale27/Go-social.git/internal/store"
	"log"
	"math/rand"
)

var usernames = []string{
	"astro_wolf", "byte_blaze", "ninja_coder", "silent_arrow", "tech_surge",
	"code_mancer", "skyline_dev", "glitch_rider", "zero_trace", "pixel_wave",
	"quantum_jumper", "neon_venom", "cryptic_blade", "echo_shadow", "debug_warrior",
	"phantom_root", "dark_loop", "matrix_dreamer", "terminal_fury", "cyber_hawk",
	"ghost_packet", "hex_rider", "data_pulse", "script_sniper", "react_raven",
	"storm_logic", "binary_ninja", "core_phantom", "shadow_loop", "bit_crusher",
	"warp_falcon", "code_lancer", "silent_pixel", "ajax_frost", "null_strike",
	"giga_spike", "fusion_blitz", "bot_exile", "snip3r_soul", "cmd_mystic",
	"dev_stalker", "code_raptor", "vault_hacker", "prime_loop", "arcane_stack",
	"stack_sprinter", "root_void", "jolt_raven", "static_rush", "overclock_x",
}

var postTitles = []string{
	"Getting Started with Go: A Beginner's Guide",
	"Understanding Goroutines and Concurrency in Go",
	"10 Go Tips Every Developer Should Know",
	"How to Build REST APIs with Go",
	"Mastering the Go Module System",
	"Go Interfaces Explained with Examples",
	"Working with JSON in Go",
	"Building a CLI Tool in Go",
	"Understanding Pointers in Go",
	"Error Handling Best Practices in Go",
	"Structs and Methods in Golang",
	"Writing Unit Tests in Go",
	"Gin vs Echo: Choosing a Web Framework in Go",
	"Using Channels to Communicate Between Goroutines",
	"Database Integration in Go Using GORM",
	"Creating a Simple Web Server with net/http",
	"Concurrency Patterns in Golang",
	"Go vs Rust: A Developer's Perspective",
	"How the Go Scheduler Works",
	"Deploying Go Applications with Docker",
}

var postContents = []string{
	"This post walks you through installing Go, setting up your workspace, and writing your first Go program.",
	"Learn how to run tasks concurrently using goroutines and manage them effectively with channels.",
	"A collection of practical tips and tricks that can level up your Go development experience.",
	"Step-by-step guide on building a RESTful API using Go’s standard library and best practices.",
	"Understand how Go modules work and how to manage dependencies efficiently.",
	"A beginner-friendly breakdown of interfaces and how they support polymorphism in Go.",
	"This post covers encoding and decoding JSON in Go, along with tips for working with structs.",
	"Learn how to build a simple and efficient command-line tool in Go using the flag and cobra packages.",
	"Explore the concept of pointers in Go, how they differ from other languages, and how to use them safely.",
	"Learn about Go’s error handling approach and how to write clean, readable error-handling code.",
	"Deep dive into structs, embedding, methods, and how they help organize code in Go.",
	"This post introduces Go’s built-in testing package and shows how to write and run unit tests.",
	"Compare two popular Go web frameworks, Gin and Echo, based on performance, ease of use, and features.",
	"Understand how channels work and how they enable safe communication between goroutines.",
	"A tutorial on integrating relational databases with Go using the GORM ORM.",
	"Learn how to build a simple web server using the standard net/http package in Go.",
	"Explore common concurrency patterns such as fan-out/fan-in, worker pools, and pipelines in Go.",
	"A balanced comparison of Go and Rust for systems programming, performance, and learning curve.",
	"This post demystifies how Go schedules goroutines and what makes its concurrency model efficient.",
	"Learn how to containerize and deploy Go applications using Docker and Dockerfiles.",
}

var postTags = []string{
	"go", "beginner", "setup", "getting-started",
	"concurrency", "goroutines", "channels",
	"tips", "best-practices", "productivity",
	"api", "rest", "web", "http",
	"modules", "dependency-management",
	"interfaces", "oop", "abstraction",
	"json", "encoding", "decoding",
	"cli", "command-line", "tools",
	"pointers", "memory", "basics",
	"error-handling", "clean-code",
	"structs", "methods",
	"testing", "unit-tests", "tdd",
	"web-frameworks", "gin", "echo",
	"database", "gorm", "orm",
	"patterns", "fan-out", "fan-in", "worker-pool",
	"rust", "comparison", "performance",
	"scheduler", "runtime",
	"docker", "deployment", "containers",
}

var postComments = []string{
	"Awesome intro to Go! Helped me get started quickly.",
	"The concurrency explanation was super clear, thanks!",
	"Great tips — I didn’t know about the `go vet` tool!",
	"This API tutorial was exactly what I needed!",
	"Modules were confusing before, this post clarified a lot.",
	"Interfaces finally make sense now, thanks to your examples.",
	"Worked perfectly for my JSON parsing task. Thanks!",
	"Just built my first CLI tool in Go — loved this!",
	"I finally understand pointers after reading this. Great job!",
	"Very helpful overview of error handling in Go.",
	"Clear and concise explanation of structs and methods!",
	"Unit testing in Go is simpler than I expected. Nice post!",
	"Helpful comparison — ended up going with Gin!",
	"Channels were tricky, but your visuals really helped.",
	"Was stuck on GORM setup — this post saved me!",
	"Loved the simplicity of using net/http to build a server.",
	"The concurrency patterns section was a gem!",
	"Really enjoyed the Go vs Rust breakdown — insightful!",
	"Scheduler internals were a mystery to me before this!",
	"Thanks for the Docker guide — deployed my first Go app!",
}

func Seed(store store2.Storage) {
	ctx := context.Background()
	users := generateUsers(100)

	for _, user := range users {

		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user", user, err)
			return
		}

	}

	posts := generatePosts(200, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post", post, err)
			return
		}
	}

	comments := generateComments(50, users, posts)
	for _, comment := range comments {
		postid := comment.PostID
		if _, err := store.Comments.create(ctx, postid); err != nil {
			log.Println("Error getting post", postid, err)
		}
	}

	return
}

func generateUsers(num int) []*store2.User {

	/*
		A slice is a collection of elements.
		In this case, the elements are pointers to store2.User objects.
		Instead of storing the actual store2.User objects in the slice,
		the slice stores the memory addresses (pointers) of those objects.
	*/

	users := make([]*store2.User, num)

	for i := 0; i < num; i++ {

		users[i] = &store2.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}

	}
	return users
}

func generatePosts(num int, users []*store2.User) []*store2.Post {
	posts := make([]*store2.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))] // selects a random user

		posts[i] = &store2.Post{
			UserID:  user.Id,
			Title:   postTitles[rand.Intn(len(postTitles))],
			Content: postContents[rand.Intn(len(postContents))],
			Tags: []string{
				postTags[rand.Intn(len(postTags))],
				postTags[rand.Intn(len(postTags))],
			},
		}
	}
	return posts
}

func generateComments(num int, users []*store2.User, posts []*store2.Post) []*store2.Comment {
	comments := make([]*store2.Comment, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]
		comments[i] = &store2.Comment{
			PostID:  post.Id,
			UserID:  user.Id,
			Content: postContents[rand.Intn(len(postContents))],
		}
	}
}
