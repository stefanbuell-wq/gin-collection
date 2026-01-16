package admin

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// ServerHandler handles server management operations
type ServerHandler struct {
	projectPath string
}

// NewServerHandler creates a new server handler
func NewServerHandler(projectPath string) *ServerHandler {
	return &ServerHandler{
		projectPath: projectPath,
	}
}

// ContainerStatus represents a Docker container's status
type ContainerStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Health  string `json:"health"`
	Ports   string `json:"ports"`
	Created string `json:"created"`
}

// ServerStatus represents the overall server status
type ServerStatus struct {
	Containers  []ContainerStatus `json:"containers"`
	DiskUsage   string            `json:"disk_usage"`
	MemoryUsage string            `json:"memory_usage"`
	Uptime      string            `json:"uptime"`
	LastDeploy  string            `json:"last_deploy,omitempty"`
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

// GetStatus handles GET /admin/api/v1/server/status
func (h *ServerHandler) GetStatus(c *gin.Context) {
	status := ServerStatus{
		Containers: h.getContainerStatus(),
		DiskUsage:  h.getDiskUsage(),
		MemoryUsage: h.getMemoryUsage(),
		Uptime:     h.getUptime(),
	}

	c.JSON(200, status)
}

// getContainerStatus gets Docker container status
func (h *ServerHandler) getContainerStatus() []ContainerStatus {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}|{{.Status}}|{{.Ports}}")
	output, err := cmd.Output()
	if err != nil {
		logger.Error("Failed to get container status", "error", err.Error())
		return nil
	}

	var containers []ContainerStatus
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 2 {
			status := ContainerStatus{
				Name:   parts[0],
				Status: parts[1],
			}
			if len(parts) >= 3 {
				status.Ports = parts[2]
			}
			// Extract health from status if present
			if strings.Contains(status.Status, "(healthy)") {
				status.Health = "healthy"
			} else if strings.Contains(status.Status, "(unhealthy)") {
				status.Health = "unhealthy"
			} else if strings.Contains(status.Status, "Up") {
				status.Health = "running"
			} else {
				status.Health = "stopped"
			}
			containers = append(containers, status)
		}
	}

	return containers
}

// getDiskUsage gets disk usage
func (h *ServerHandler) getDiskUsage() string {
	cmd := exec.Command("df", "-h", "/")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) >= 2 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 5 {
			return fmt.Sprintf("%s used of %s (%s)", fields[2], fields[1], fields[4])
		}
	}
	return "unknown"
}

// getMemoryUsage gets memory usage
func (h *ServerHandler) getMemoryUsage() string {
	cmd := exec.Command("free", "-h")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) >= 2 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 3 {
			return fmt.Sprintf("%s used of %s", fields[2], fields[1])
		}
	}
	return "unknown"
}

// getUptime gets system uptime
func (h *ServerHandler) getUptime() string {
	cmd := exec.Command("uptime", "-p")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

// DeployRequest represents a deployment request
type DeployRequest struct {
	Services []string `json:"services"` // e.g., ["api", "frontend", "admin-frontend"]
	Pull     bool     `json:"pull"`     // Whether to git pull first
	NoCache  bool     `json:"no_cache"` // Whether to build with --no-cache
}

// Deploy handles POST /admin/api/v1/server/deploy
func (h *ServerHandler) Deploy(c *gin.Context) {
	var req DeployRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var outputs []string

	// Git pull if requested
	if req.Pull {
		logger.Info("Executing git pull")
		cmd := exec.Command("git", "pull")
		cmd.Dir = h.projectPath
		output, err := cmd.CombinedOutput()
		outputs = append(outputs, fmt.Sprintf("Git Pull:\n%s", string(output)))
		if err != nil {
			c.JSON(500, CommandResult{
				Success: false,
				Output:  strings.Join(outputs, "\n\n"),
				Error:   err.Error(),
			})
			return
		}
	}

	// Build and restart services
	for _, service := range req.Services {
		logger.Info("Deploying service", "service", service)

		// Build
		buildArgs := []string{"compose", "build"}
		if req.NoCache {
			buildArgs = append(buildArgs, "--no-cache")
		}
		buildArgs = append(buildArgs, service)

		cmd := exec.Command("docker", buildArgs...)
		cmd.Dir = h.projectPath
		output, err := cmd.CombinedOutput()
		outputs = append(outputs, fmt.Sprintf("Build %s:\n%s", service, string(output)))
		if err != nil {
			c.JSON(500, CommandResult{
				Success: false,
				Output:  strings.Join(outputs, "\n\n"),
				Error:   fmt.Sprintf("Failed to build %s: %s", service, err.Error()),
			})
			return
		}

		// Up
		cmd = exec.Command("docker", "compose", "up", "-d", service)
		cmd.Dir = h.projectPath
		output, err = cmd.CombinedOutput()
		outputs = append(outputs, fmt.Sprintf("Up %s:\n%s", service, string(output)))
		if err != nil {
			c.JSON(500, CommandResult{
				Success: false,
				Output:  strings.Join(outputs, "\n\n"),
				Error:   fmt.Sprintf("Failed to start %s: %s", service, err.Error()),
			})
			return
		}
	}

	c.JSON(200, CommandResult{
		Success: true,
		Output:  strings.Join(outputs, "\n\n"),
	})
}

// RestartService handles POST /admin/api/v1/server/restart/:service
func (h *ServerHandler) RestartService(c *gin.Context) {
	service := c.Param("service")
	if service == "" {
		c.JSON(400, gin.H{"error": "Service name required"})
		return
	}

	// Validate service name (security)
	validServices := map[string]bool{
		"api":            true,
		"frontend":       true,
		"admin-frontend": true,
		"mysql":          true,
		"redis":          true,
		"nginx":          true,
	}

	if !validServices[service] {
		c.JSON(400, gin.H{"error": "Invalid service name"})
		return
	}

	logger.Info("Restarting service", "service", service)

	var cmd *exec.Cmd
	if service == "nginx" {
		cmd = exec.Command("sudo", "systemctl", "reload", "nginx")
	} else {
		cmd = exec.Command("docker", "compose", "restart", service)
		cmd.Dir = h.projectPath
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, CommandResult{
			Success: false,
			Output:  string(output),
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, CommandResult{
		Success: true,
		Output:  string(output),
	})
}

// GetLogs handles GET /admin/api/v1/server/logs/:service
func (h *ServerHandler) GetLogs(c *gin.Context) {
	service := c.Param("service")
	if service == "" {
		c.JSON(400, gin.H{"error": "Service name required"})
		return
	}

	lines := c.DefaultQuery("lines", "100")

	// Validate service name (security)
	validServices := map[string]bool{
		"api":            true,
		"frontend":       true,
		"admin-frontend": true,
		"mysql":          true,
		"redis":          true,
	}

	if !validServices[service] {
		c.JSON(400, gin.H{"error": "Invalid service name"})
		return
	}

	// Map service names to container names
	containerNames := map[string]string{
		"api":            "gin-collection-api",
		"frontend":       "gin-collection-frontend",
		"admin-frontend": "gin-collection-admin",
		"mysql":          "gin-collection-mysql",
		"redis":          "gin-collection-redis",
	}

	containerName := containerNames[service]
	cmd := exec.Command("docker", "logs", "--tail", lines, containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, CommandResult{
			Success: false,
			Output:  string(output),
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, CommandResult{
		Success: true,
		Output:  string(output),
	})
}

// StreamLogs handles GET /admin/api/v1/server/logs/:service/stream (SSE)
func (h *ServerHandler) StreamLogs(c *gin.Context) {
	service := c.Param("service")
	if service == "" {
		c.JSON(400, gin.H{"error": "Service name required"})
		return
	}

	// Validate service name
	containerNames := map[string]string{
		"api":            "gin-collection-api",
		"frontend":       "gin-collection-frontend",
		"admin-frontend": "gin-collection-admin",
		"mysql":          "gin-collection-mysql",
		"redis":          "gin-collection-redis",
	}

	containerName, ok := containerNames[service]
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid service name"})
		return
	}

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	cmd := exec.Command("docker", "logs", "-f", "--tail", "50", containerName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.SSEvent("error", err.Error())
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.SSEvent("error", err.Error())
		return
	}

	if err := cmd.Start(); err != nil {
		c.SSEvent("error", err.Error())
		return
	}

	// Stream stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			c.SSEvent("log", scanner.Text())
			c.Writer.Flush()
		}
	}()

	// Stream stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			c.SSEvent("log", scanner.Text())
			c.Writer.Flush()
		}
	}()

	// Keep connection alive until client disconnects
	<-c.Request.Context().Done()
	cmd.Process.Kill()
}

// NginxReload handles POST /admin/api/v1/server/nginx/reload
func (h *ServerHandler) NginxReload(c *gin.Context) {
	logger.Info("Reloading nginx configuration")

	// Test config first
	testCmd := exec.Command("sudo", "nginx", "-t")
	testOutput, err := testCmd.CombinedOutput()
	if err != nil {
		c.JSON(500, CommandResult{
			Success: false,
			Output:  string(testOutput),
			Error:   "Nginx config test failed",
		})
		return
	}

	// Reload
	reloadCmd := exec.Command("sudo", "systemctl", "reload", "nginx")
	reloadOutput, err := reloadCmd.CombinedOutput()
	if err != nil {
		c.JSON(500, CommandResult{
			Success: false,
			Output:  string(reloadOutput),
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, CommandResult{
		Success: true,
		Output:  "Nginx reloaded successfully\n" + string(testOutput),
	})
}

// QuickActions represents available quick actions
type QuickAction struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Dangerous   bool   `json:"dangerous"`
}

// GetQuickActions handles GET /admin/api/v1/server/actions
func (h *ServerHandler) GetQuickActions(c *gin.Context) {
	actions := []QuickAction{
		{
			ID:          "pull-all",
			Name:        "Git Pull & Deploy All",
			Description: "Pull latest code and rebuild all services",
			Icon:        "rocket",
			Dangerous:   false,
		},
		{
			ID:          "deploy-api",
			Name:        "Deploy API",
			Description: "Rebuild and restart API service",
			Icon:        "server",
			Dangerous:   false,
		},
		{
			ID:          "deploy-frontend",
			Name:        "Deploy Frontend",
			Description: "Rebuild and restart frontend",
			Icon:        "layout",
			Dangerous:   false,
		},
		{
			ID:          "deploy-admin",
			Name:        "Deploy Admin Panel",
			Description: "Rebuild and restart admin frontend",
			Icon:        "shield",
			Dangerous:   false,
		},
		{
			ID:          "reload-nginx",
			Name:        "Reload Nginx",
			Description: "Test and reload nginx configuration",
			Icon:        "refresh-cw",
			Dangerous:   false,
		},
		{
			ID:          "restart-all",
			Name:        "Restart All Containers",
			Description: "Restart all Docker containers",
			Icon:        "power",
			Dangerous:   true,
		},
	}

	c.JSON(200, gin.H{"actions": actions})
}

// ExecuteAction handles POST /admin/api/v1/server/actions/:action
func (h *ServerHandler) ExecuteAction(c *gin.Context) {
	action := c.Param("action")

	logger.Info("Executing quick action", "action", action)

	switch action {
	case "pull-all":
		h.executePullAll(c)
	case "deploy-api":
		h.executeDeploy(c, "api")
	case "deploy-frontend":
		h.executeDeploy(c, "frontend")
	case "deploy-admin":
		h.executeDeploy(c, "admin-frontend")
	case "reload-nginx":
		h.NginxReload(c)
	case "restart-all":
		h.executeRestartAll(c)
	default:
		c.JSON(400, gin.H{"error": "Unknown action"})
	}
}

func (h *ServerHandler) executePullAll(c *gin.Context) {
	var outputs []string

	// Git pull
	cmd := exec.Command("git", "pull")
	cmd.Dir = h.projectPath
	output, err := cmd.CombinedOutput()
	outputs = append(outputs, fmt.Sprintf("=== Git Pull ===\n%s", string(output)))
	if err != nil {
		c.JSON(500, CommandResult{Success: false, Output: strings.Join(outputs, "\n"), Error: err.Error()})
		return
	}

	// Build all services
	services := []string{"api", "frontend", "admin-frontend"}
	for _, svc := range services {
		cmd = exec.Command("docker", "compose", "build", svc)
		cmd.Dir = h.projectPath
		output, _ = cmd.CombinedOutput()
		outputs = append(outputs, fmt.Sprintf("=== Build %s ===\n%s", svc, string(output)))
	}

	// Restart all
	cmd = exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = h.projectPath
	output, err = cmd.CombinedOutput()
	outputs = append(outputs, fmt.Sprintf("=== Start Services ===\n%s", string(output)))

	if err != nil {
		c.JSON(500, CommandResult{Success: false, Output: strings.Join(outputs, "\n"), Error: err.Error()})
		return
	}

	c.JSON(200, CommandResult{Success: true, Output: strings.Join(outputs, "\n")})
}

func (h *ServerHandler) executeDeploy(c *gin.Context, service string) {
	var outputs []string
	startTime := time.Now()

	// Build
	cmd := exec.Command("docker", "compose", "build", "--no-cache", service)
	cmd.Dir = h.projectPath
	output, err := cmd.CombinedOutput()
	outputs = append(outputs, fmt.Sprintf("=== Build %s ===\n%s", service, string(output)))
	if err != nil {
		c.JSON(500, CommandResult{Success: false, Output: strings.Join(outputs, "\n"), Error: err.Error()})
		return
	}

	// Up
	cmd = exec.Command("docker", "compose", "up", "-d", service)
	cmd.Dir = h.projectPath
	output, err = cmd.CombinedOutput()
	outputs = append(outputs, fmt.Sprintf("=== Start %s ===\n%s", service, string(output)))
	if err != nil {
		c.JSON(500, CommandResult{Success: false, Output: strings.Join(outputs, "\n"), Error: err.Error()})
		return
	}

	duration := time.Since(startTime).Round(time.Second)
	outputs = append(outputs, fmt.Sprintf("\nDeployment completed in %s", duration))

	c.JSON(200, CommandResult{Success: true, Output: strings.Join(outputs, "\n")})
}

func (h *ServerHandler) executeRestartAll(c *gin.Context) {
	cmd := exec.Command("docker", "compose", "restart")
	cmd.Dir = h.projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		c.JSON(500, CommandResult{Success: false, Output: string(output), Error: err.Error()})
		return
	}

	c.JSON(200, CommandResult{Success: true, Output: string(output)})
}
