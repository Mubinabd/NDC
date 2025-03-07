package handlers

import (
	"context"
	"google.golang.org/protobuf/encoding/protojson"
	"strconv"

	pb "posts/internal/pkg/genproto"

	"github.com/gin-gonic/gin"
)

// CreateLog creates a new log
// @Summary Create Log
// @Description Create a new log
// @Tags Log
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param log body pb.LogCreateRequest true "Log data"
// @Success 200 {string} string pb.LogCreateResponse "Log created successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/logs/create [post]
func (h *Handler) CreateLog(c *gin.Context) {
	var req pb.LogCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.ERROR.Println("Failed to bind request:", err)
		c.JSON(400, "Invalid request: "+err.Error())
		return
	}

	res, err := h.Clients.Log.Create(c, &req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to create log:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, res)
}

// GetLogDetail retrieves a log by ID
// @Summary Get Log
// @Description Retrieve a log by its ID
// @Tags Log
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Log ID"
// @Success 200 {object} pb.LogGetResponse "Log details"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Log not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/logs/{id} [get]
func (h *Handler) GetLogDetail(c *gin.Context) {
	logID := c.Param("id")

	id, err := strconv.ParseInt(logID, 10, 64)
	if err != nil {
		h.Logger.ERROR.Println("Invalid log ID:", err)
		c.JSON(400, "Invalid log ID")
		return
	}

	req := &pb.GetId{Id: id}
	res, err := h.Clients.Log.GetDetail(c, req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to get log:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, res)
}

// UpdateLog updates an existing log
// @Summary Update Log
// @Description Update an existing log
// @Tags Log
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param log body pb.LogUpdateRequest true "Log data"
// @Success 200 {string} string "Log updated successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/logs/update [patch]
func (h *Handler) UpdateLog(c *gin.Context) {
	var body pb.LogUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.Logger.ERROR.Println("Failed to bind request:", err)
		c.JSON(400, "Invalid request: "+err.Error())
		return
	}

	req := pb.LogUpdateRequest{
		Id:          body.Id,
		Level:       body.Level,
		Message:     body.Message,
		ServiceName: body.ServiceName,
	}

	//_, err := h.Clients.Log.Update(c, &req)
	//if err != nil {
	//	h.Logger.ERROR.Println("Failed to update log:", err)
	//	c.JSON(500, "Internal server error: "+err.Error())
	//	return
	//}

	input, err := protojson.Marshal(&req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to marshal req:", err)
		c.JSON(500, "Invalid request: "+err.Error())
		return
	}
	if err := h.Producer.ProduceMessages("log-update", input); err != nil {
		h.Logger.ERROR.Println("Failed to produce message:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, "Log updated successfully")
}

// DeleteLog deletes a log by ID
// @Summary Delete Log
// @Description Delete a log by its ID
// @Tags Log
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Log ID"
// @Success 200 {string} string "Log deleted successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Log not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/logs/{id} [delete]
func (h *Handler) DeleteLog(c *gin.Context) {
	logID := c.Param("id")

	id, err := strconv.ParseInt(logID, 10, 64)
	if err != nil {
		h.Logger.ERROR.Println("Invalid log ID:", err)
		c.JSON(400, "Invalid log ID")
		return
	}

	req := &pb.GetId{Id: id}

	//_, err := h.Clients.Log.Delete(c, req)
	//if err != nil {
	//	h.Logger.ERROR.Println("Failed to delete log:", err)
	//	c.JSON(500, "Internal server error: "+err.Error())
	//	return
	//}

	input, err := protojson.Marshal(req)
	if err != nil {
		h.Logger.ERROR.Println("Failed to marshal req:", err)
		c.JSON(500, "Invalid request: "+err.Error())
		return
	}
	if err := h.Producer.ProduceMessages("log-delete", input); err != nil {
		h.Logger.ERROR.Println("Failed to produce message:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, "Log deleted successfully")
}

// GetLogList retrieves a list of logs with pagination
// @Summary Get Logs
// @Description Retrieve a list of logs with pagination
// @Tags Log
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param page query int false "Page"
// @Param level query string false "Level"
// @Param service_name query string false "ServiceName"
// @Success 200 {object} pb.LogGetAll "List of logs"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/logs/list [get]
func (h *Handler) GetLogList(c *gin.Context) {
	var filter pb.FilterLog

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

	if level := c.Query("level"); level != "" {
		filter.Level = level
	}

	if servicename := c.Query("servicename"); servicename != "" {
		filter.ServiceName = servicename
	}

	res, err := h.Clients.Log.GetList(context.Background(), &filter)
	if err != nil {
		h.Logger.ERROR.Println("Failed to list logs:", err)
		c.JSON(500, "Internal server error: "+err.Error())
		return
	}

	c.JSON(200, res)
}
