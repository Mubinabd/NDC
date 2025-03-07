package handlers

import (
	"context"
	"google.golang.org/protobuf/encoding/protojson"
	"strconv"

	pb "posts/internal/pkg/genproto"

	"github.com/gin-gonic/gin"
)

// Create creates a new post
// @Summary Create Post
// @Description Create a new post
// @Tags Post
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post body pb.PostCreateRequest true "Post data"
// @Success 200 {string} string pb.PostCreateResponse "Post created successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/posts/create [post]
func (h *Handler) Create(c *gin.Context) {
	var req pb.PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.ERROR.Println("Failed to bind request:", err)
		c.JSON(400, "Invalid request: "+err.Error())
		return
	}

	res, err := h.Clients.Post.Create(c, &req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to create post:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, res)
}

// GetDetail retrieves a post by ID
// @Summary Get Post
// @Description Retrieve a post by its ID
// @Tags Post
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} pb.PostGetResponse "Post details"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Post not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/posts/{id} [get]
func (h *Handler) GetDetail(c *gin.Context) {
	postID := c.Param("id")

	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		h.Logger.ERROR.Println("Invalid log ID:", err)
		c.JSON(400, "Invalid log ID")
		return
	}

	req := &pb.GetById{Id: id}
	res, err := h.Clients.Post.GetDetail(c, req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to get post:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, res)
}

// UpdatePost updates an existing post
// @Summary Update Post
// @Description Update an existing post
// @Tags Post
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post body pb.PostUpdateRequest true "Post data"
// @Success 200 {string} string "Post updated successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/posts/update [patch]
func (h *Handler) UpdatePost(c *gin.Context) {
	var body pb.PostUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.Logger.ERROR.Println("Failed to bind request:", err)
		c.JSON(400, "Invalid request: "+err.Error())
		return
	}

	req := pb.PostUpdateRequest{
		Id:      body.Id,
		UserId:  body.UserId,
		Title:   body.Title,
		Content: body.Content,
	}

	//_, err := h.Clients.Post.Update(c, &req)
	//if err != nil {
	//	h.Logger.ERROR.Println("Failed to update post:", err)
	//	c.JSON(500, "Internal server error: "+err.Error())
	//	return
	//}

	input, err := protojson.Marshal(&req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to marshal req:", err)
		c.JSON(500, "Invalid request: "+err.Error())
		return
	}
	if err := h.Producer.ProduceMessages("post-update", input); err != nil {
		h.Logger.ERROR.Println("Failed to produce message:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, "Post updated successfully")
}

// Delete deletes a post by ID
// @Summary Delete Post
// @Description Delete a post by its ID
// @Tags Post
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {string} string "Post deleted successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Post not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/posts/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	postID := c.Param("id")

	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		h.Logger.ERROR.Println("Invalid log ID:", err)
		c.JSON(400, "Invalid log ID")
		return
	}

	req := &pb.GetById{Id: id}

	//_, err := h.Clients.Post.Delete(c, req)
	//if err != nil {
	//	h.Logger.ERROR.Println("Failed to delete post:", err)
	//	c.JSON(500, "Internal server error: "+err.Error())
	//	return
	//}

	input, err := protojson.Marshal(req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to marshal req:", err)
		c.JSON(500, "Invalid request: "+err.Error())
		return
	}
	if err := h.Producer.ProduceMessages("post-delete", input); err != nil {
		h.Logger.ERROR.Println("Failed to produce message:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, "Post deleted successfully")
}

// GetList retrieves a list of posts with pagination
// @Summary Get Posts
// @Description Retrieve a list of posts with pagination
// @Tags Post
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param page query int false "Page"
// @Param user_id query int false "UserID"
// @Param content query string false "Content"
// @Success 200 {object} pb.PostGetAll "List of posts"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/posts/list [get]
func (h *Handler) GetList(c *gin.Context) {
	var filter pb.FilterPost

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = int64(l)
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = int64(o)
		}
	}
	if user_id := c.Query("user_id"); user_id != "" {
		if o, err := strconv.Atoi(user_id); err == nil {
			filter.UserId = int64(o)
		}
	}
	if content := c.Query("content"); content != "" {
		filter.Content = content
	}

	res, err := h.Clients.Post.GetList(context.Background(), &filter)
	if err != nil {
		h.Logger.ERROR.Println("Failed to list posts:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, res)
}
