package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootDir string

func main() {
	// For now, we'll use the environment variable ARCH_VPS_ROOT or default to /mnt/storage/arch-vps-server
	rootDir = os.Getenv("ARCH_VPS_ROOT")
	if rootDir == "" {
		rootDir = "/mnt/storage/arch-vps-server"
	}

	var rootCmd = &cobra.Command{
		Use:   "lab",
		Short: "Arch VPS Lab Manager",
		Long:  `A CLI tool to manage your Arch VPS server lab, including proxy and projects.`,
	}

	rootCmd.AddCommand(proxyCmd())
	rootCmd.AddCommand(networkCmd())
	rootCmd.AddCommand(projectCmd())
	rootCmd.AddCommand(statusCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Proxy Commands
func proxyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Manage the Caddy proxy",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Start the proxy",
		Run: func(cmd *cobra.Command, args []string) {
			err := runCommand(filepath.Join(rootDir, "proxy"), "docker", "compose", "up", "-d")
			if err != nil {
				fmt.Printf("Error starting proxy: %v\n", err)
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Stop the proxy",
		Run: func(cmd *cobra.Command, args []string) {
			err := runCommand(filepath.Join(rootDir, "proxy"), "docker", "compose", "down")
			if err != nil {
				fmt.Printf("Error stopping proxy: %v\n", err)
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "restart",
		Short: "Restart the proxy",
		Run: func(cmd *cobra.Command, args []string) {
			proxyDir := filepath.Join(rootDir, "proxy")
			runCommand(proxyDir, "docker", "compose", "down")
			err := runCommand(proxyDir, "docker", "compose", "up", "-d")
			if err != nil {
				fmt.Printf("Error restarting proxy: %v\n", err)
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "logs",
		Short: "Follow proxy logs",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(filepath.Join(rootDir, "proxy"), "docker", "compose", "logs", "-f")
		},
	})

	return cmd
}

// Network Commands
func networkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "network",
		Short: "Manage Docker networks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Creating arch-vps-net network...")
			err := exec.Command("docker", "network", "create", "arch-vps-net").Run()
			if err != nil {
				if strings.Contains(err.Error(), "already exists") {
					fmt.Println("Network arch-vps-net already exists.")
				} else {
					fmt.Printf("Error creating network: %v\n", err)
				}
			} else {
				fmt.Println("Network arch-vps-net created successfully.")
			}
		},
	}
}

// Project Commands
func projectCmd() *cobra.Command {
	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
	}

	projectCmd.AddCommand(&cobra.Command{
		Use:   "update [project-path]",
		Short: "Update a project (git pull & docker compose up)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectPath := args[0]
			if !filepath.IsAbs(projectPath) {
				// If it doesn't exist relative to CWD, try relative to projects/
				if _, err := os.Stat(projectPath); os.IsNotExist(err) {
					projectPath = filepath.Join(rootDir, "projects", projectPath)
				}
			}

			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				fmt.Printf("Error: project directory does not exist: %s\n", projectPath)
				return
			}

			fmt.Printf("Updating project in %s...\n", projectPath)
			if err := runCommand(projectPath, "git", "pull"); err != nil {
				fmt.Printf("Error pulling git: %v\n", err)
			}

			composeFile := "docker-compose.yml"
			if _, err := os.Stat(filepath.Join(projectPath, "compose.yml")); err == nil {
				composeFile = "compose.yml"
			}

			if err := runCommand(projectPath, "docker", "compose", "-f", composeFile, "up", "-d", "--build"); err != nil {
				fmt.Printf("Error starting project: %v\n", err)
			}
		},
	})

	// Scaffolding command
	var domain string
	var port int
	addCmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add and scaffold a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			projectPath := filepath.Join(rootDir, "projects", name)

			if _, err := os.Stat(projectPath); err == nil {
				fmt.Printf("Error: project %s already exists at %s\n", name, projectPath)
				return
			}

			// 1. Create directory
			os.MkdirAll(projectPath, 0755)

			// 2. Generate docker-compose.yml
			composeContent := fmt.Sprintf(`services:
  %s:
    build: .
    container_name: %s
    restart: unless-stopped
    expose:
      - "%d"
    networks:
      - arch-vps-net

networks:
  arch-vps-net:
    external: true
`, name, name, port)
			os.WriteFile(filepath.Join(projectPath, "docker-compose.yml"), []byte(composeContent), 0644)

			// 3. Add to Caddyfile
			caddyPath := filepath.Join(rootDir, "proxy", "Caddyfile")
			caddyEntry := fmt.Sprintf("\n%s {\n    reverse_proxy %s:%d\n}\n", domain, name, port)
			
			f, err := os.OpenFile(caddyPath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Error opening Caddyfile: %v\n", err)
			} else {
				defer f.Close()
				f.WriteString(caddyEntry)
				fmt.Println("Added entry to Caddyfile.")
			}

			fmt.Printf("Project %s scaffolded at %s\n", name, projectPath)
			fmt.Println("Next steps:")
			fmt.Printf("1. Add your source code to %s\n", projectPath)
			fmt.Println("2. Run: lab proxy restart")
		},
	}
	addCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain for the project")
	addCmd.Flags().IntVarP(&port, "port", "p", 80, "Internal port of the container")
	addCmd.MarkFlagRequired("domain")
	projectCmd.AddCommand(addCmd)

	return projectCmd
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show status of all lab containers",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(rootDir, "docker", "ps", "--filter", "network=arch-vps-net")
		},
	}
}
