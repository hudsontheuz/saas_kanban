package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	"github.com/hudsontheuz/saas_kanban/internal/shared/envx"
	task "github.com/hudsontheuz/saas_kanban/internal/task/domain"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
	team "github.com/hudsontheuz/saas_kanban/internal/team/domain"
	teamrepo "github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/gorm/repo"
	user "github.com/hudsontheuz/saas_kanban/internal/user/domain"
	userrepo "github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/gorm/repo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	if err := envx.Load(); err != nil {
		t.Fatalf("erro carregando .env: %v", err)
	}

	dsn := strings.TrimSpace(os.Getenv("TEST_DB_URL"))
	if dsn == "" {
		dsn = strings.TrimSpace(os.Getenv("DB_URL"))
	}
	if dsn == "" {
		t.Skip("defina TEST_DB_URL ou DB_URL para rodar os testes com GORM")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("erro abrindo banco de teste: %v", err)
	}

	applySchema(t, db)
	resetTables(t, db)
	return db
}

func applySchema(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlPath := filepath.Join("..", "migrations", "001_start_system.up.sql")
	content, err := os.ReadFile(sqlPath)
	if err != nil {
		t.Fatalf("erro lendo migration: %v", err)
	}

	for _, stmt := range splitSQLStatements(string(content)) {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("erro aplicando migration em %q: %v", stmt, err)
		}
	}
}

func splitSQLStatements(sqlText string) []string {
	var cleaned []string
	for _, line := range strings.Split(sqlText, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		cleaned = append(cleaned, line)
	}

	parts := strings.Split(strings.Join(cleaned, "\n"), ";")
	stmts := make([]string, 0, len(parts))
	for _, part := range parts {
		stmt := strings.TrimSpace(part)
		if stmt == "" || strings.EqualFold(stmt, "BEGIN") || strings.EqualFold(stmt, "COMMIT") {
			continue
		}
		stmts = append(stmts, stmt)
	}
	return stmts
}

func resetTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	stmt := `TRUNCATE TABLE tarefa, projeto, equipe_membro, equipe, usuario RESTART IDENTITY CASCADE;`
	if err := db.Exec(stmt).Error; err != nil {
		t.Fatalf("erro limpando tabelas de teste: %v", err)
	}
}

func seedUser(t *testing.T, db *gorm.DB, nome string) user.UserID {
	t.Helper()

	repo := userrepo.NewUserRepo(db)
	slug := strings.ToLower(strings.ReplaceAll(nome, " ", "_"))

	u, err := user.NovoUsuario(
		nome,
		fmt.Sprintf("%s_%d@test.local", slug, time.Now().UnixNano()),
		"hash-teste",
	)
	if err != nil {
		t.Fatalf("erro criando user de teste: %v", err)
	}

	if err := repo.Salvar(u); err != nil {
		t.Fatalf("erro salvando user de teste: %v", err)
	}

	return u.ID()
}

func seedTeam(t *testing.T, db *gorm.DB, nome string, leaderID user.UserID) team.TeamID {
	t.Helper()

	repo := teamrepo.NewTeamRepo(db)

	tm, err := team.NovaTeam(nome, leaderID)
	if err != nil {
		t.Fatalf("erro criando team de teste: %v", err)
	}

	if err := repo.Salvar(tm); err != nil {
		t.Fatalf("erro salvando team de teste: %v", err)
	}

	return tm.ID()
}

func seedProject(
	t *testing.T,
	db *gorm.DB,
	teamID team.TeamID,
	nome string,
	cfg project.ConfiguracoesProject,
) project.ProjectID {
	t.Helper()

	repo := projectrepo.NewProjectRepo(db)

	p, err := project.NovoProject(teamID, nome, cfg)
	if err != nil {
		t.Fatalf("erro criando project de teste: %v", err)
	}

	if err := repo.Salvar(p); err != nil {
		t.Fatalf("erro salvando project de teste: %v", err)
	}

	return p.ID()
}

func seedTask(t *testing.T, db *gorm.DB, projectID project.ProjectID, titulo string) task.TaskID {
	t.Helper()

	repo := taskrepo.NewTaskRepo(db)

	tk, err := task.NovaTask(projectID, titulo)
	if err != nil {
		t.Fatalf("erro criando task de teste: %v", err)
	}

	if err := repo.Salvar(tk); err != nil {
		t.Fatalf("erro salvando task de teste: %v", err)
	}

	return tk.ID()
}

func mustInt64ID(t *testing.T, id string) int64 {
	t.Helper()

	n, err := strconv.ParseInt(strings.TrimSpace(id), 10, 64)
	if err != nil {
		t.Fatalf("id inválido %q: %v", id, err)
	}

	return n
}
