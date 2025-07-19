package unit

import (
	"testing"

	"github.com/oiahoon/termonaut/internal/privacy"
)

func TestPasswordSanitization(t *testing.T) {
	sanitizer := privacy.NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "mysql password",
			input:    "mysql -u root -p'secret123' -h localhost",
			expected: "mysql -u root -p'[REDACTED]' -h localhost",
		},
		{
			name:     "postgres password",
			input:    "psql postgresql://user:password@localhost/db",
			expected: "psql postgresql://user:[REDACTED]@localhost/db",
		},
		{
			name:     "ssh with password",
			input:    "sshpass -p 'mypassword' ssh user@host",
			expected: "sshpass -p '[REDACTED]' ssh user@host",
		},
		{
			name:     "curl with auth",
			input:    "curl -u admin:secret123 https://api.example.com",
			expected: "curl -u admin:[REDACTED] https://api.example.com",
		},
		{
			name:     "environment variable password",
			input:    "export DB_PASSWORD=supersecret",
			expected: "export DB_PASSWORD=[REDACTED]",
		},
		{
			name:     "no password",
			input:    "ls -la /home/user",
			expected: "ls -la /home/user",
		},
		{
			name:     "git clone with token",
			input:    "git clone https://token123@github.com/user/repo.git",
			expected: "git clone https://[REDACTED]@github.com/user/repo.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeCommand(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s\nGot: %s", tt.expected, result)
			}
		})
	}
}

func TestTokenSanitization(t *testing.T) {
	sanitizer := privacy.NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "github token",
			input:    "export GITHUB_TOKEN=ghp_1234567890abcdef",
			expected: "export GITHUB_TOKEN=[REDACTED]",
		},
		{
			name:     "api key",
			input:    "curl -H 'Authorization: Bearer sk-1234567890abcdef' https://api.openai.com",
			expected: "curl -H 'Authorization: Bearer [REDACTED]' https://api.openai.com",
		},
		{
			name:     "aws credentials",
			input:    "aws s3 ls --access-key AKIAIOSFODNN7EXAMPLE",
			expected: "aws s3 ls --access-key [REDACTED]",
		},
		{
			name:     "jwt token",
			input:    "curl -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c'",
			expected: "curl -H 'Authorization: [REDACTED]'",
		},
		{
			name:     "docker login",
			input:    "docker login -u user -p dckr_pat_1234567890abcdef",
			expected: "docker login -u user -p [REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeCommand(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s\nGot: %s", tt.expected, result)
			}
		})
	}
}

func TestURLSanitization(t *testing.T) {
	sanitizer := privacy.NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "http url with credentials",
			input:    "curl https://user:pass@example.com/api",
			expected: "curl https://[REDACTED]@example.com/api",
		},
		{
			name:     "ftp url with credentials",
			input:    "wget ftp://admin:secret@ftp.example.com/file.txt",
			expected: "wget ftp://[REDACTED]@ftp.example.com/file.txt",
		},
		{
			name:     "database connection string",
			input:    "psql postgres://user:password@localhost:5432/mydb",
			expected: "psql postgres://user:[REDACTED]@localhost:5432/mydb",
		},
		{
			name:     "redis url",
			input:    "redis-cli -u redis://user:password@localhost:6379/0",
			expected: "redis-cli -u redis://user:[REDACTED]@localhost:6379/0",
		},
		{
			name:     "clean url",
			input:    "curl https://api.example.com/users",
			expected: "curl https://api.example.com/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeCommand(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s\nGot: %s", tt.expected, result)
			}
		})
	}
}

func TestFilePathSanitization(t *testing.T) {
	sanitizer := privacy.NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "home directory path",
			input:    "cat /Users/john/Documents/secret.txt",
			expected: "cat /Users/[USER]/Documents/secret.txt",
		},
		{
			name:     "generic home path",
			input:    "ls /home/alice/projects",
			expected: "ls /home/[USER]/projects",
		},
		{
			name:     "system path",
			input:    "sudo cat /etc/passwd",
			expected: "sudo cat /etc/passwd",
		},
		{
			name:     "relative path",
			input:    "cat ./config/database.yml",
			expected: "cat ./config/database.yml",
		},
		{
			name:     "temp directory",
			input:    "ls /tmp/user_session_12345",
			expected: "ls /tmp/user_session_12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeCommand(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s\nGot: %s", tt.expected, result)
			}
		})
	}
}

func TestEmailSanitization(t *testing.T) {
	sanitizer := privacy.NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "git config email",
			input:    "git config user.email john.doe@company.com",
			expected: "git config user.email [EMAIL]",
		},
		{
			name:     "sendmail command",
			input:    "echo 'test' | mail -s 'subject' user@example.com",
			expected: "echo 'test' | mail -s 'subject' [EMAIL]",
		},
		{
			name:     "curl with email parameter",
			input:    "curl -d 'email=admin@site.com' https://api.example.com",
			expected: "curl -d 'email=[EMAIL]' https://api.example.com",
		},
		{
			name:     "no email",
			input:    "echo 'hello world'",
			expected: "echo 'hello world'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeCommand(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s\nGot: %s", tt.expected, result)
			}
		})
	}
}

func TestIPAddressSanitization(t *testing.T) {
	sanitizer := privacy.NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ssh to ip",
			input:    "ssh user@192.168.1.100",
			expected: "ssh user@[IP]",
		},
		{
			name:     "ping ip",
			input:    "ping 10.0.0.1",
			expected: "ping [IP]",
		},
		{
			name:     "curl to ip",
			input:    "curl http://172.16.0.50:8080/api",
			expected: "curl http://[IP]:8080/api",
		},
		{
			name:     "localhost",
			input:    "curl http://localhost:3000",
			expected: "curl http://localhost:3000",
		},
		{
			name:     "127.0.0.1",
			input:    "telnet 127.0.0.1 80",
			expected: "telnet 127.0.0.1 80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeCommand(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s\nGot: %s", tt.expected, result)
			}
		})
	}
}

func TestSanitizerConfiguration(t *testing.T) {
	// Test with different sanitizer configurations
	tests := []struct {
		name           string
		enablePassword bool
		enableURL      bool
		enableEmail    bool
		enableIP       bool
		input          string
		expected       string
	}{
		{
			name:           "all enabled",
			enablePassword: true,
			enableURL:      true,
			enableEmail:    true,
			enableIP:       true,
			input:          "mysql -u root -p'secret' user@example.com 192.168.1.1",
			expected:       "mysql -u root -p'[REDACTED]' [EMAIL] [IP]",
		},
		{
			name:           "only password",
			enablePassword: true,
			enableURL:      false,
			enableEmail:    false,
			enableIP:       false,
			input:          "mysql -u root -p'secret' user@example.com 192.168.1.1",
			expected:       "mysql -u root -p'[REDACTED]' user@example.com 192.168.1.1",
		},
		{
			name:           "none enabled",
			enablePassword: false,
			enableURL:      false,
			enableEmail:    false,
			enableIP:       false,
			input:          "mysql -u root -p'secret' user@example.com 192.168.1.1",
			expected:       "mysql -u root -p'secret' user@example.com 192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := privacy.SanitizerConfig{
				SanitizePasswords: tt.enablePassword,
				SanitizeURLs:      tt.enableURL,
				SanitizeEmails:    tt.enableEmail,
				SanitizeIPs:       tt.enableIP,
			}
			sanitizer := privacy.NewSanitizerWithConfig(config)
			
			result := sanitizer.SanitizeCommand(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s\nGot: %s", tt.expected, result)
			}
		})
	}
}

func TestSanitizerPerformance(t *testing.T) {
	sanitizer := privacy.NewSanitizer()
	
	// Test with a long command
	longCommand := "mysql -u root -p'verylongpasswordwithmanycharacters' -h database.example.com -P 3306 -D production_database --execute=\"SELECT * FROM users WHERE email='user@example.com' AND created_at > '2023-01-01'\" --batch --raw"
	
	// Run sanitization multiple times to test performance
	for i := 0; i < 1000; i++ {
		result := sanitizer.SanitizeCommand(longCommand)
		if result == longCommand {
			t.Error("Command should have been sanitized")
		}
	}
}

func TestEdgeCases(t *testing.T) {
	sanitizer := privacy.NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty command",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only",
			input:    "   \t\n  ",
			expected: "   \t\n  ",
		},
		{
			name:     "special characters",
			input:    "echo 'password=secret123' | grep -v '^#'",
			expected: "echo 'password=[REDACTED]' | grep -v '^#'",
		},
		{
			name:     "multiple passwords",
			input:    "script.sh --db-pass secret1 --api-key secret2",
			expected: "script.sh --db-pass [REDACTED] --api-key [REDACTED]",
		},
		{
			name:     "unicode characters",
			input:    "echo 'пароль=секрет123'",
			expected: "echo 'пароль=[REDACTED]'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeCommand(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s\nGot: %s", tt.expected, result)
			}
		})
	}
}
