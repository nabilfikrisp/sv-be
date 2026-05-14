// Package v1 provides post REST API handlers.
package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nabilfikrisp/sv-be/internal/controller/restapi/v1/request"
	"github.com/nabilfikrisp/sv-be/internal/dto"
	"github.com/nabilfikrisp/sv-be/internal/entity"
)

// @Summary      Create a new post
// @Description  Create a post with the provided details
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        request body request.CreatePost true "Create Post Body"
// @Success      201  {object}  entity.Post
// @Failure      400  {object}  response.Error  "Invalid request or validation error"
// @Failure      500  {object}  response.Error  "Internal server error"
// @Router       /article [post]
func (r *V1) createPost(c *gin.Context) {
	var body request.CreatePost

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - createPost")
		validationErrorResponse(c, err)
		return
	}

	post, err := r.uc_post.Create(c.Request.Context(), dto.PostCreate{
		Title:    body.Title,
		Content:  body.Content,
		Category: body.Category,
		Status:   body.Status,
	})
	if err != nil {
		r.l.Error(err, "restapi - v1 - createPost")
		if errors.Is(err, entity.ErrPostTitleEmpty) {
			errorResponse(c, http.StatusBadRequest, "Title cannot be empty")
			return
		}
		if errors.Is(err, entity.ErrPostStatusInvalid) {
			errorResponse(c, http.StatusBadRequest, "Invalid status")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	c.JSON(http.StatusCreated, post)
}

// @Summary      Get a post by ID
// @Description  Retrieve details of a single post using its unique ID
// @Tags         posts
// @Produce      json
// @Param        id   path      int  true  "Post ID"
// @Success      200  {object}  entity.Post
// @Failure      400  {object}  response.Error  "Invalid ID"
// @Failure      404  {object}  response.Error  "Post not found"
// @Failure      500  {object}  response.Error  "Internal server error"
// @Router       /article/{id} [get]
func (r *V1) getPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	post, err := r.uc_post.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrPostNotFound) {
			errorResponse(c, http.StatusNotFound, "Post not found")
			return
		}

		r.l.Error(err, "restapi - v1 - getPostByID")
		errorResponse(c, http.StatusInternalServerError, "Failed to get post")
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary      List posts
// @Description  Retrieve a paginated list of posts with optional filtering
// @Tags         posts
// @Produce      json
// @Param        status   query     string  false  "Filter by status (publish, draft, thrash)"
// @Param        limit    query     int     false  "Page limit (default 10)"
// @Param        offset   query     int     false  "Page offset (default 0)"
// @Success      200  {object}  map[string]interface{} "Returns {posts: []entity.Post, total: int}"
// @Failure      400  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /article [get]
func (r *V1) listPosts(c *gin.Context) {
	filter := request.PostFilter{
		Limit:  new(uint64(10)),
		Offset: new(uint64(0)),
	}

	if v := c.Query("status"); v != "" {
		status := entity.PostStatus(v)
		if !status.Valid() {
			errorResponse(c, http.StatusBadRequest, "invalid status")
			return
		}
		filter.Status = &status
	}

	if v, err := strconv.ParseUint(c.Query("limit"), 10, 64); err == nil {
		filter.Limit = &v
	}
	if v, err := strconv.ParseUint(c.Query("offset"), 10, 64); err == nil {
		filter.Offset = &v
	}

	posts, total, err := r.uc_post.List(c.Request.Context(), filter.ToDTO())
	if err != nil {
		r.l.Error(err, "restapi - v1 - listPosts")
		if errors.Is(err, entity.ErrPostStatusInvalid) {
			errorResponse(c, http.StatusBadRequest, "Invalid status filter")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to list posts")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
	})
}

// @Summary      Update a post
// @Description  Update existing post details by ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id      path      int                 true  "Post ID"
// @Param        request body      request.UpdatePost  true  "Update Post Body"
// @Success      200  {object}  entity.Post
// @Failure      400  {object}  response.Error
// @Failure      404  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /article/{id} [patch]
func (r *V1) updatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var body request.UpdatePost

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - updatePost")
		validationErrorResponse(c, err)
		return
	}

	post, err := r.uc_post.Update(c.Request.Context(), id, dto.PostUpdate{
		Title:    body.Title,
		Content:  body.Content,
		Category: body.Category,
		Status:   body.Status,
	})
	if err != nil {
		if errors.Is(err, entity.ErrPostNotFound) {
			errorResponse(c, http.StatusNotFound, "Post not found")
			return
		}

		if errors.Is(err, entity.ErrPostStatusInvalid) {
			errorResponse(c, http.StatusBadRequest, "Invalid status")
			return
		}

		r.l.Error(err, "restapi - v1 - updatePost")
		errorResponse(c, http.StatusInternalServerError, "Failed to update post")
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary      Delete a post
// @Description  Remove a post from the system by ID
// @Tags         posts
// @Param        id   path      int  true  "Post ID"
// @Success      204  "No Content"
// @Failure      404  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /article/{id} [delete]
func (r *V1) deletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := r.uc_post.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, entity.ErrPostNotFound) {
			errorResponse(c, http.StatusNotFound, "Post not found")
			return
		}

		r.l.Error(err, "restapi - v1 - deletePost")
		errorResponse(c, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	c.Status(http.StatusNoContent)
}
