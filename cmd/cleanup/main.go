package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nais/dependencytrack/pkg/client"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	adminPwd := os.Getenv("ADMIN_PASSWORD")
	baseUrl := os.Getenv("BASE_URL")
	deleteProjects := os.Getenv("DELETE_PROJECTS")
	deleteTeams := os.Getenv("DELETE_TEAMS")
	deleteOidcUsers := os.Getenv("DELETE_OIDC_USERS")

	c := client.New(baseUrl, "admin", adminPwd)

	ctx := context.Background()
	if deleteProjects != "" {
		projects, err := c.GetProjects(ctx)
		if err != nil {
			panic(err)
		}

		var projectlog string
		log.Infof("Projects to delete: %d", len(projects))
		for _, project := range projects {
			if err := c.DeleteProject(ctx, project.Uuid); err != nil {
				log.Errorf("Error deleting project %s: %s", project.Name, err)
			}
			projectlog += fmt.Sprintf("Deleted project %s\n", project.Name)
		}
		log.Info(projectlog)

		err = c.PortfolioRefresh(ctx)
		if err != nil {
			log.Errorf("Error refreshing portfolio: %s", err)
		}
	}

	if deleteTeams != "" {
		teams, err := c.GetTeams(context.Background())
		if err != nil {
			panic(err)
		}
		var teamlog string
		log.Infof("Teams to delete: %d", len(teams))
		for _, team := range teams {
			if team.Name != "Administrators" && team.Name != "Automation" && team.Name != "Portfolio Managers" {
				if err := c.DeleteTeam(ctx, team.Uuid); err != nil {
					log.Fatalf("Error deleting team %s: %s", team.Name, err)
				}
				teamlog += fmt.Sprintf("Deleted team %s\n", team.Name)
			}
		}
		log.Info(teamlog)
	}

	if deleteOidcUsers != "" {
		users, err := c.GetOidcUsers(context.Background())
		if err != nil {
			panic(err)
		}
		var userlog string
		log.Infof("Users to delete: %d", len(users))
		for _, user := range users {
			if err := c.DeleteOidcUser(ctx, user.Username); err != nil {
				log.Fatalf("Error deleting user %s: %s", user.Username, err)
			}
			userlog += fmt.Sprintf("Deleted user %s\n", user.Username)
		}
		log.Info(userlog)
	}
}
