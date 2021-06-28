package main

import (
	"fmt"
	pb "grpccrud/proto"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewDataserviceClient(conn)
	g := gin.Default()
	grp1 := g.Group("/api/v1")
	{
		grp1.GET("user/:a/:b", func(ctx *gin.Context) {
			a := ctx.Param("a")

			b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
				return
			}

			req := &pb.ReadRequest{Api: a, Id: int64(b)}
			if response, err := client.Read(ctx, req); err == nil {
				ctx.JSON(http.StatusOK, gin.H{
					"result": fmt.Sprint(response.Api, response.ToDo),
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		})
		// grp1.POST("user/:a", func(ctx *gin.Context) {
		// 	a := ctx.Param("a")

		// 	b := ctx.BindJSON(&pb.ToDo)

		// 	req := &pb.ReadRequest{Api: a, Id: int64(b)}
		// 	if response, err := client.Read(ctx, req); err == nil {
		// 		ctx.JSON(http.StatusOK, gin.H{
		// 			"result": fmt.Sprint(response.Result),
		// 		})
		// 	} else {
		// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	}
		// })
		//grp1.GET("user/:id", Controllers.GetUserByID)
		// grp1.PUT("user/:id", Controllers.UpdateUser)
		// grp1.DELETE("user/:id", Controllers.DeleteUser)
	}
	// g.GET("/add/:a/:b", func(ctx *gin.Context) {
	// 	a, err := strconv.ParseUint(ctx.Param("a"), 10, 64)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
	// 		return
	// 	}

	// 	b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
	// 		return
	// 	}

	// 	req := &proto.Request{A: int64(a), B: int64(b)}
	// 	if response, err := client.Add(ctx, req); err == nil {
	// 		ctx.JSON(http.StatusOK, gin.H{
	// 			"result": fmt.Sprint(response.Result),
	// 		})
	// 	} else {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	}
	// })

	// g.GET("/mult/:a/:b", func(ctx *gin.Context) {
	// 	a, err := strconv.ParseUint(ctx.Param("a"), 10, 64)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
	// 		return
	// 	}
	// 	b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
	// 		return
	// 	}
	// 	req := &proto.Request{A: int64(a), B: int64(b)}

	// 	if response, err := client.Multiply(ctx, req); err == nil {
	// 		ctx.JSON(http.StatusOK, gin.H{
	// 			"result": fmt.Sprint(response.Result),
	// 		})
	// 	} else {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	}
	// })

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
