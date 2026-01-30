package email

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client *resend.Client
	appUrl string
}

func NewEmailService(apiKey string, appUrl string) *EmailService {
	client := resend.NewClient(apiKey)
	return &EmailService{
		client: client,
		appUrl: appUrl,
	}
}

func (s *EmailService) SendInvitation(toEmail, orgName, inviterName, token string) error {
	inviteLink := fmt.Sprintf("%s/invite/accept?token=%s", s.appUrl, token)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
			<h2>You've been invited to join %s</h2>
			<p>%s has invited you to join <strong>%s</strong> on Chis.</p>
			<p>Click the button below to accept the invitation:</p>
			<div style="text-align: center; margin: 30px 0;">
				<a href="%s"
				   style="background-color: #4F46E5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; font-weight: bold;">
					Accept Invitation
				</a>
			</div>
			<p style="color: #666; font-size: 14px;">This invitation will expire in 7 days.</p>
			<p style="color: #666; font-size: 14px;">If you didn't expect this invitation, you can ignore this email.</p>
		</div>
	`, orgName, inviterName, orgName, inviteLink)

	params := &resend.SendEmailRequest{
		From:    "Chis <noreply@trychis.com>",
		To:      []string{toEmail},
		Subject: fmt.Sprintf("You've been invited to join %s", orgName),
		Html:    htmlBody,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send invitation email: %w", err)
	}

	return nil
}
