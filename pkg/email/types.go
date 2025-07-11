package email

// TranslationSet contains all translated strings for a specific language
type TranslationSet struct {
	Title            string
	HeaderText       string
	Greeting         string
	MainText         string
	ButtonText       string
	LinkInstructions string
	SecurityNote     string
	ClosingText      string
	TeamSignature    string
	RightsText       string
	AddressText      string
}

// PasswordResetData contains all the data needed for the email template
type PasswordResetData struct {
	UserName     string
	CompanyName  string
	LogoURL      string
	ResetToken   string
	ResetURL     string
	BaseURL      string
	CurrentYear  int
	Language     string
	Translations TranslationSet
}

// InvitationData contains all the data needed for the invitation email template
type InvitationData struct {
	RecipientName string
	CompanyName   string
	LogoURL       string
	TeamName      string
	Role          string
	InviteLink    string
	CurrentYear   int
	Language      string
	Translations  InvitationTranslationSet
}

// InvitationTranslationSet contains all translated strings for invitation emails
type InvitationTranslationSet struct {
	Title            string
	HeaderText       string
	Greeting         string
	MainText         string
	AsRole           string
	ButtonText       string
	LinkInstructions string
	SecurityNote     string
	ClosingText      string
	TeamSignature    string
	RightsText       string
	AddressText      string
}

// GetTranslations returns the translations for a specific language
func GetTranslations(lang string) TranslationSet {
	translations := map[string]TranslationSet{
		"en": {
			Title:            "Reset Your Password",
			HeaderText:       "Password Reset",
			Greeting:         "Hello",
			MainText:         "We received a request to reset your password. Click the button below to create a new password. This link will expire in 24 hours.",
			ButtonText:       "Reset Password",
			LinkInstructions: "If the button doesn't work, copy and paste this link into your browser:",
			SecurityNote:     "If you didn't request a password reset, please ignore this email or contact support if you have concerns about your account security.",
			ClosingText:      "Thank you for using our service.",
			TeamSignature:    "The Neploy Team",
			RightsText:       "All rights reserved.",
			AddressText:      "123 Tech Street, Cloud City, CC 12345",
		},
		"es": {
			Title:            "Restablece tu Contraseña",
			HeaderText:       "Restablecimiento de Contraseña",
			Greeting:         "Hola",
			MainText:         "Recibimos una solicitud para restablecer tu contraseña. Haz clic en el botón a continuación para crear una nueva contraseña. Este enlace caducará en 24 horas.",
			ButtonText:       "Restablecer Contraseña",
			LinkInstructions: "Si el botón no funciona, copia y pega este enlace en tu navegador:",
			SecurityNote:     "Si no solicitaste un restablecimiento de contraseña, ignora este correo electrónico o contacta con soporte si tienes preocupaciones sobre la seguridad de tu cuenta.",
			ClosingText:      "Gracias por utilizar nuestro servicio.",
			TeamSignature:    "El Equipo de Neploy",
			RightsText:       "Todos los derechos reservados.",
			AddressText:      "Calle Tecnología 123, Ciudad Nube, CN 12345",
		},
		"pt": {
			Title:            "Redefinir Sua Senha",
			HeaderText:       "Redefinição de Senha",
			Greeting:         "Olá",
			MainText:         "Recebemos uma solicitação para redefinir sua senha. Clique no botão abaixo para criar uma nova senha. Este link expirará em 24 horas.",
			ButtonText:       "Redefinir Senha",
			LinkInstructions: "Se o botão não funcionar, copie e cole este link no seu navegador:",
			SecurityNote:     "Se você não solicitou uma redefinição de senha, ignore este e-mail ou entre em contato com o suporte se tiver preocupações sobre a segurança da sua conta.",
			ClosingText:      "Obrigado por usar nosso serviço.",
			TeamSignature:    "A Equipe Neploy",
			RightsText:       "Todos os direitos reservados.",
			AddressText:      "Rua da Tecnologia 123, Cidade Nuvem, CN 12345",
		},
		"fr": {
			Title:            "Réinitialisez Votre Mot de Passe",
			HeaderText:       "Réinitialisation du Mot de Passe",
			Greeting:         "Bonjour",
			MainText:         "Nous avons reçu une demande de réinitialisation de votre mot de passe. Cliquez sur le bouton ci-dessous pour créer un nouveau mot de passe. Ce lien expirera dans 24 heures.",
			ButtonText:       "Réinitialiser le Mot de Passe",
			LinkInstructions: "Si le bouton ne fonctionne pas, copiez et collez ce lien dans votre navigateur :",
			SecurityNote:     "Si vous n'avez pas demandé de réinitialisation de mot de passe, veuillez ignorer cet e-mail ou contacter le support si vous avez des inquiétudes concernant la sécurité de votre compte.",
			ClosingText:      "Merci d'utiliser notre service.",
			TeamSignature:    "L'équipe Neploy",
			RightsText:       "Tous droits réservés.",
			AddressText:      "123 Rue de la Technologie, Ville Nuage, VN 12345",
		},
		"zh": {
			Title:            "重置您的密码",
			HeaderText:       "密码重置",
			Greeting:         "您好",
			MainText:         "我们收到了重置您密码的请求。点击下面的按钮创建新密码。此链接将在24小时后过期。",
			ButtonText:       "重置密码",
			LinkInstructions: "如果按钮不起作用，请复制并粘贴此链接到您的浏览器：",
			SecurityNote:     "如果您没有请求重置密码，请忽略此电子邮件，或者如果您担心帐户安全，请联系支持团队。",
			ClosingText:      "感谢您使用我们的服务。",
			TeamSignature:    "Neploy团队",
			RightsText:       "保留所有权利。",
			AddressText:      "科技街123号，云城，12345",
		},
	}

	// Default to English if language not supported
	if _, ok := translations[lang]; !ok {
		lang = "en"
	}

	return translations[lang]
}

// GetInvitationTranslations returns the translations for invitation emails in a specific language
func GetInvitationTranslations(lang string) InvitationTranslationSet {
	translations := map[string]InvitationTranslationSet{
		"en": {
			Title:            "Invitation to join Neploy",
			HeaderText:       "You've Been Invited!",
			Greeting:         "Hello",
			MainText:         "You've been invited to join",
			AsRole:           "as a",
			ButtonText:       "Accept Invitation",
			LinkInstructions: "If the button doesn't work, copy and paste this link into your browser:",
			SecurityNote:     "If you didn't expect this invitation, you can safely ignore this email.",
			ClosingText:      "We look forward to having you on board.",
			TeamSignature:    "The Neploy Team",
			RightsText:       "All rights reserved.",
			AddressText:      "123 Tech Street, Cloud City, CC 12345",
		},
		"es": {
			Title:            "Invitación para unirse a Neploy",
			HeaderText:       "¡Has sido invitado!",
			Greeting:         "Hola",
			MainText:         "Has sido invitado a unirte a",
			AsRole:           "como",
			ButtonText:       "Aceptar Invitación",
			LinkInstructions: "Si el botón no funciona, copia y pega este enlace en tu navegador:",
			SecurityNote:     "Si no esperabas esta invitación, puedes ignorar este correo electrónico.",
			ClosingText:      "Esperamos contar contigo.",
			TeamSignature:    "El Equipo de Neploy",
			RightsText:       "Todos los derechos reservados.",
			AddressText:      "Calle Tecnología 123, Ciudad Nube, CN 12345",
		},
		"pt": {
			Title:            "Convite para se juntar ao Neploy",
			HeaderText:       "Você foi convidado!",
			Greeting:         "Olá",
			MainText:         "Você foi convidado para se juntar a",
			AsRole:           "como",
			ButtonText:       "Aceitar Convite",
			LinkInstructions: "Se o botão não funcionar, copie e cole este link no seu navegador:",
			SecurityNote:     "Se você não esperava este convite, pode ignorar este e-mail.",
			ClosingText:      "Esperamos tê-lo a bordo.",
			TeamSignature:    "A Equipe Neploy",
			RightsText:       "Todos os direitos reservados.",
			AddressText:      "Rua da Tecnologia 123, Cidade Nuvem, CN 12345",
		},
		"fr": {
			Title:            "Invitation à rejoindre Neploy",
			HeaderText:       "Vous avez été invité !",
			Greeting:         "Bonjour",
			MainText:         "Vous avez été invité à rejoindre",
			AsRole:           "en tant que",
			ButtonText:       "Accepter l'invitation",
			LinkInstructions: "Si le bouton ne fonctionne pas, copiez et collez ce lien dans votre navigateur :",
			SecurityNote:     "Si vous n'attendiez pas cette invitation, vous pouvez ignorer cet e-mail.",
			ClosingText:      "Nous nous réjouissons de vous avoir à bord.",
			TeamSignature:    "L'équipe Neploy",
			RightsText:       "Tous droits réservés.",
			AddressText:      "123 Rue de la Technologie, Ville Nuage, VN 12345",
		},
		"zh": {
			Title:            "邀请您加入Neploy",
			HeaderText:       "您收到了邀请！",
			Greeting:         "您好",
			MainText:         "您被邀请加入",
			AsRole:           "作为",
			ButtonText:       "接受邀请",
			LinkInstructions: "如果按钮不起作用，请复制并粘贴此链接到您的浏览器：",
			SecurityNote:     "如果您没有预料到这个邀请，您可以安全地忽略这封电子邮件。",
			ClosingText:      "我们期待您的加入。",
			TeamSignature:    "Neploy团队",
			RightsText:       "保留所有权利。",
			AddressText:      "科技街123号，云城，12345",
		},
	}

	// Default to English if language not supported
	if _, ok := translations[lang]; !ok {
		lang = "en"
	}

	return translations[lang]
}
