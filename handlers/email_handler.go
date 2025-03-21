package handlers

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/rogerjeasy/go-letusconnect/config"
)

// sendEmail is a helper function to send an email with a specified subject and body
func sendEmail(toEmail, subject, body string) error {
	auth := smtp.PlainAuth("", config.SenderEmail, config.SenderPass, config.SMTPHost)

	// Email message format
	msg := fmt.Sprintf("From: %s <%s>\r\n", config.SenderName, config.SenderEmail) +
		fmt.Sprintf("To: %s\r\n", toEmail) +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort),
		auth,
		config.SenderEmail,
		[]string{toEmail},
		[]byte(msg),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// SendAutomaticEmail sends a professional thank-you email to the sender
func SendAutomaticEmail(toEmail, senderName string) error {
	subject := "Thank You for Reaching Out to LetUsConnect!"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: Arial, sans-serif;
      color: #333333;
      background-color: #f9f9f9;
      padding: 20px;
      text-align: center;
    }
    .container {
      max-width: 600px;
      margin: 0 auto;
      background: #ffffff;
      padding: 30px;
      border-radius: 10px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    h2 {
      color: #4A90E2;
    }
    p {
      font-size: 16px;
      line-height: 1.6;
    }
    .footer {
      margin-top: 30px;
      font-size: 14px;
      color: #777777;
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>Thank You, %s!</h2>
    <p>
      We appreciate you reaching out to <strong>LetUsConnect Support</strong>. Your message has been received, and our team will get back to you as soon as possible.
    </p>
    <p>
      We are here to assist you with any questions, concerns, or feedback you may have.
    </p>
    <p>
      In the meantime, feel free to explore our <a href="https://letusconnect.vercel.app/help" style="color: #4A90E2;">Help Center</a> for useful resources and FAQs.
    </p>
    <p class="footer">
      Best regards, <br>
      <strong>The LetUsConnect Team</strong>
    </p>
  </div>
</body>
</html>`, senderName)

	return sendEmail(toEmail, subject, body)
}

// SendProjectJoinNotification sends an email to the user when they join a project
func SendProjectJoinNotification(toEmail, userName, projectName string) error {
	subject := "Welcome to the Project - Let's Collaborate!"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: Arial, sans-serif;
      color: #333333;
      background-color: #f9f9f9;
      padding: 20px;
      text-align: center;
    }
    .container {
      max-width: 600px;
      margin: 0 auto;
      background: #ffffff;
      padding: 30px;
      border-radius: 10px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    h2 {
      color: #4A90E2;
    }
    p {
      font-size: 16px;
      line-height: 1.6;
    }
    .footer {
      margin-top: 30px;
      font-size: 14px;
      color: #777777;
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>Congratulations, %s! 🎉</h2>
    <p>
      You have successfully joined the project: <strong>"%s"</strong>.
    </p>
    <p>
      We are excited to have you on board and look forward to seeing your contributions. Collaborate with your team, share ideas, and make a positive impact!
    </p>
    <p>
      To get started, visit your project dashboard and connect with your fellow team members.
    </p>
    <p class="footer">
      Best regards, <br>
      <strong>The LetUsConnect Team</strong>
    </p>
  </div>
</body>
</html>`, userName, projectName)

	return sendEmail(toEmail, subject, body)
}

// SendJoinRequestAcceptedEmail sends an email to the user when their join request is accepted
func SendJoinRequestAcceptedEmail(toEmail, userName, projectName string) error {
	subject := "Your Project Join Request Was Accepted!"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: Arial, sans-serif;
      color: #333333;
      background-color: #f9f9f9;
      padding: 20px;
      text-align: center;
    }
    .container {
      max-width: 600px;
      margin: 0 auto;
      background: #ffffff;
      padding: 30px;
      border-radius: 10px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    h2 {
      color: #4A90E2;
    }
    p {
      font-size: 16px;
      line-height: 1.6;
    }
    .footer {
      margin-top: 30px;
      font-size: 14px;
      color: #777777;
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>Good News, %s! 🎉</h2>
    <p>
      Your request to join the project <strong>"%s"</strong> has been accepted!
    </p>
    <p>
      You are now part of the project team. Connect with your team members, collaborate on tasks, and let's make great things happen together.
    </p>
    <p>
      Visit the project dashboard to get started.
    </p>
    <p class="footer">
      Best regards, <br>
      <strong>The LetUsConnect Team</strong>
    </p>
  </div>
</body>
</html>`, userName, projectName)

	return sendEmail(toEmail, subject, body)
}

// SendJoinRequestRejectedEmail sends an email to the user when their join request is rejected
func SendJoinRequestRejectedEmail(toEmail, userName, projectName string) error {
	subject := "Your Project Join Request Was Rejected"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: Arial, sans-serif;
      color: #333333;
      background-color: #f9f9f9;
      padding: 20px;
      text-align: center;
    }
    .container {
      max-width: 600px;
      margin: 0 auto;
      background: #ffffff;
      padding: 30px;
      border-radius: 10px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    h2 {
      color: #E74C3C;
    }
    p {
      font-size: 16px;
      line-height: 1.6;
    }
    .highlight {
      color: #E67E22;
      font-weight: bold;
    }
    .footer {
      margin-top: 30px;
      font-size: 14px;
      color: #777777;
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>Dear %s,</h2>
    <p>
      Unfortunately, your request to join the project <strong>"%s"</strong> was not accepted. 😔
    </p>
    <p>
      We understand this may be disappointing, but don't be discouraged! There are many other exciting projects waiting for your participation on <span class="highlight">LetUsConnect</span>.
    </p>
    <p>
      Explore new projects and find the perfect opportunity to contribute your skills and talents! 🚀
    </p>
    <p>
      We hope to see you collaborating with other teams soon.
    </p>
    <p class="footer">
      Best regards, <br>
      <strong>The LetUsConnect Team</strong>
    </p>
  </div>
</body>
</html>`, userName, projectName)

	return sendEmail(toEmail, subject, body)
}

// SendWelcomeEmail sends a welcome email to the newly registered user
func SendWelcomeEmail(toEmail, username, logoURL string) error {
	auth := smtp.PlainAuth("", config.SenderEmail, config.SenderPass, config.SMTPHost)

	// Email subject and body with the embedded logo (if available)
	subject := "Welcome to LetUsConnect - Your Journey Starts Here!"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: Arial, sans-serif;
      color: #333333;
      background-color: #f9f9f9;
      padding: 20px;
      text-align: center;
    }
    .container {
      max-width: 600px;
      margin: 0 auto;
      background: #ffffff;
      padding: 30px;
      border-radius: 10px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    h2 {
      color: #4A90E2;
    }
    p {
      font-size: 16px;
      line-height: 1.6;
    }
    .features {
      text-align: left;
      margin-top: 20px;
    }
    .features li {
      margin-bottom: 10px;
    }
    .button {
      display: inline-block;
      background-color: #4A90E2;
      color: #ffffff;
      padding: 12px 24px;
      margin: 20px 0;
      text-decoration: none;
      border-radius: 5px;
      font-weight: bold;
    }
    .logo {
      max-width: 150px;
      margin: 20px 0;
    }
  </style>
</head>
<body>
  <div class="container">
    %s
    <h2>Welcome to LetUsConnect, %s! 🎉</h2>
    <p>
      We are thrilled to have you join our growing community of students, alumni, and professionals.
      LetUsConnect is your gateway to networking, collaboration, and career growth.
    </p>

    <p>
      Here’s what you can do with LetUsConnect:
    </p>

    <ul class="features">
      <li><strong>🤝 Connect:</strong> Meet fellow students, alumni, and professionals in your field.</li>
      <li><strong>🚀 Collaborate:</strong> Join exciting projects or create your own to work together.</li>
      <li><strong>📈 Grow:</strong> Discover new learning opportunities and career paths.</li>
      <li><strong>🗂️ Showcase:</strong> Share your skills and accomplishments with the community.</li>
    </ul>

    <p>
      LetUsConnect is designed to help you build meaningful relationships and unlock new opportunities. We believe that together, we can achieve more!
    </p>

    <a href="https://letusconnect.vercel.app/get-started" class="button">Get Started Now</a>

    <p>
      Need more information? Visit our <a href="https://letusconnect.vercel.app/help" style="color: #4A90E2;">Help Center</a> or check out our <a href="https://letusconnect.vercel.app/about-us" style="color: #4A90E2;">About Page</a>.
    </p>

    <p>
      If you have any questions, don’t hesitate to reach out to our support team at <a href="mailto:support@letusconnect.com" style="color: #4A90E2;">support@letusconnect.com</a> or at <a href="https://letusconnect.vercel.app/contact-us" style="color: #4A90E2;">Online Contact Form</a>.
    </p>

    <p>
      Let’s connect, collaborate, and grow together!
    </p>

    <p>
      Best regards, <br>
      <strong>The LetUsConnect Team</strong>
    </p>
  </div>
</body>
</html>`,
		// Conditionally include the logo if the URL is provided
		getLogoHTML(logoURL),
		username,
	)

	// Email message format
	msg := fmt.Sprintf("From: %s <%s>\r\n", config.SenderName, config.SenderEmail) +
		fmt.Sprintf("To: %s\r\n", toEmail) +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort),
		auth,
		config.SenderEmail,
		[]string{toEmail},
		[]byte(msg),
	)

	if err != nil {
		return fmt.Errorf("failed to send welcome email: %v", err)
	}

	return nil
}

// getLogoHTML returns the HTML for the logo if the URL is provided
func getLogoHTML(logoURL string) string {
	if strings.TrimSpace(logoURL) != "" {
		return fmt.Sprintf(`<img src="%s" alt="LetUsConnect Logo" class="logo" />`, logoURL)
	}
	return ""
}

// SendJoinRequestSubmittedEmail sends an email to the user when they request to join a project
func SendJoinRequestSubmittedEmail(toEmail, userName, projectName string) error {
	subject := "Your Join Request Has Been Received!"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: Arial, sans-serif;
      color: #333333;
      background-color: #f9f9f9;
      padding: 20px;
      text-align: center;
    }
    .container {
      max-width: 600px;
      margin: 0 auto;
      background: #ffffff;
      padding: 30px;
      border-radius: 10px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    h2 {
      color: #4A90E2;
    }
    p {
      font-size: 16px;
      line-height: 1.6;
    }
    .highlight {
      color: #E67E22;
      font-weight: bold;
    }
    .footer {
      margin-top: 30px;
      font-size: 14px;
      color: #777777;
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>Thank You, %s! 🙌</h2>
    <p>
      Your request to join the project <strong>"%s"</strong> has been successfully received.
    </p>
    <p>
      The project owner and team members will review your request shortly. We will notify you as soon as a decision has been made.
    </p>
    <p>
      In the meantime, feel free to explore other projects and opportunities on <span class="highlight">LetUsConnect</span>!
    </p>
    <p>
      Stay tuned, and we hope to see you collaborating soon! 🚀
    </p>
    <p class="footer">
      Best regards, <br>
      <strong>The LetUsConnect Team</strong>
    </p>
  </div>
</body>
</html>`, userName, projectName)

	return sendEmail(toEmail, subject, body)
}

// SendProjectInvitationEmail sends an email to a user who has been invited to join a project
func SendProjectInvitationEmail(toEmail, userName, projectName, inviterName string) error {
	subject := "You’ve Been Invited to Join a Project on LetUsConnect!"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: Arial, sans-serif;
      color: #333333;
      background-color: #f9f9f9;
      padding: 20px;
      text-align: center;
    }
    .container {
      max-width: 600px;
      margin: 0 auto;
      background: #ffffff;
      padding: 30px;
      border-radius: 10px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    h2 {
      color: #4A90E2;
    }
    p {
      font-size: 16px;
      line-height: 1.6;
    }
    .button {
      display: inline-block;
      background-color: #4A90E2;
      color: #ffffff;
      padding: 12px 24px;
      margin: 20px 0;
      text-decoration: none;
      border-radius: 5px;
      font-weight: bold;
    }
    .footer {
      margin-top: 30px;
      font-size: 14px;
      color: #777777;
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>You've Been Invited, %s! 🎉</h2>
    <p>
      <strong>%s</strong> has invited you to join the project: <strong>"%s"</strong>.
    </p>
    <p>
      This is a great opportunity to collaborate, contribute your skills, and be part of an exciting project.
    </p>
    <p>
      Click the button below to accept the invitation and join the project team.
    </p>
    <a href="https://letusconnect.vercel.app/projects/%s" class="button">Accept Invitation</a>
    <p>
      If you’re not sure or have questions, feel free to reach out to <strong>%s</strong> or visit the project page for more details.
    </p>
    <p class="footer">
      Best regards, <br>
      <strong>The LetUsConnect Team</strong>
    </p>
  </div>
</body>
</html>`, userName, inviterName, projectName, projectName, inviterName)

	return sendEmail(toEmail, subject, body)
}

// SendEmailToNotifyRequestJoinProject sends emails to notify project members about a new join request
func SendEmailToNotifyRequestJoinProject(toEmails []string, requestingUserName, requestingUserUID, projectName string) error {
	subject := "New Project Join Request Received"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<style>
  body {
    font-family: Arial, sans-serif;
    color: #333333;
    background-color: #f9f9f9;
    padding: 20px;
    text-align: center;
  }
  .container {
    max-width: 600px;
    margin: 0 auto;
    background: #ffffff;
    padding: 30px;
    border-radius: 10px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }
  h2 {
    color: #4A90E2;
  }
  p {
    font-size: 16px;
    line-height: 1.6;
  }
  .button {
    display: inline-block;
    background-color: #4A90E2;
    color: #ffffff;
    padding: 12px 24px;
    margin: 20px 0;
    text-decoration: none;
    border-radius: 5px;
    font-weight: bold;
    transition: background-color 0.3s ease;
  }
  .button:hover {
    background-color: #357ABD;
  }
  .profile-link {
    color: #4A90E2;
    text-decoration: none;
    font-weight: bold;
    border-bottom: 2px solid transparent;
    transition: border-color 0.3s ease;
  }
  .profile-link:hover {
    border-bottom-color: #4A90E2;
  }
  .footer {
    margin-top: 30px;
    font-size: 14px;
    color: #777777;
    border-top: 1px solid #eee;
    padding-top: 20px;
  }
  .highlight {
    background-color: #F8F9FA;
    border-radius: 5px;
    padding: 15px;
    margin: 15px 0;
  }
</style>
</head>
<body>
<div class="container">
  <h2>Exciting News! New Join Request 📬</h2>
  <p>
    <a href="https://letusconnect.vercel.app/profile/%s" class="profile-link">%s</a> has expressed interest in joining your project!
  </p>
  <div class="highlight">
    <p style="font-size: 18px; margin: 0;">
      Project: <strong>"%s"</strong>
    </p>
  </div>
  <p>
    This is a great opportunity to potentially welcome a new member to your team. We encourage you to:
  </p>
  <ul style="text-align: left; padding-left: 30px;">
    <li>Review their profile to learn about their skills and experience</li>
    <li>Check their portfolio and past contributions</li>
    <li>Consider how they might enhance your project</li>
  </ul>
  <p>
    Ready to review their request?
  </p>
  <a href="https://letusconnect.vercel.app/projects/manage-requests" class="button">Review Join Request →</a>
  <p style="font-size: 14px; color: #666;">
    You can easily accept or decline this request through your project management dashboard.
  </p>
  <p class="footer">
    Together, let's build amazing things! 🚀<br><br>
    Best regards,<br>
    <strong>The LetUsConnect Team</strong><br>
    <span style="font-size: 12px; color: #999;">Building Connections, Fostering Collaboration</span>
  </p>
</div>
</body>
</html>`, requestingUserUID, requestingUserName, projectName)

	// Send email to each recipient
	for _, email := range toEmails {
		if err := sendEmail(email, subject, body); err != nil {
			log.Printf("Failed to send join request notification email to %s: %v", email, err)
			continue
		}
	}

	return nil
}
