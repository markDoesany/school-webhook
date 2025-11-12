package constants

const (
	WelcomeMessage = "Welcome to School Assistant!\n\n" +
		"I'm your all-in-one academic and financial assistant. From checking grades to managing tuition, I'm here to help you stay on top of your school responsibilities."

	AccountDeactivatedMessage = "Your account has been deactivated. Kindly contact your school admin for assistance."

	AboutUsMessage = "𝗔𝗯𝗼𝘂𝘁 𝗢𝘂𝗿 𝗦𝗰𝗵𝗼𝗼𝗹 𝗔𝘀𝘀𝗶𝘀𝘁𝗮𝗻𝘁 🎓\n\n" +
		"School Assistant is your all-in-one academic and financial assistant for student life.\n\n" +
		"From checking grades to managing tuition, this chatbot helps you stay on top of your school responsibilities, right here in Messenger.\n\n" +
		"What You Can Do with School Assistant:\n" +
		"• 📚 Check your grades – View your academic performance anytime.\n• 💳 Manage school fees – Track balances, due dates, and settle payments easily.\n• 📢 View announcements – Stay informed about events, updates, and important deadlines.\n• 🤝 Get support – Need help with anything school-related? Just ask.\n\n" +
		"Need assistance? I'm always here to help you out."

	TalkToHumanMessage = "🛠️ 𝗧𝗮𝗹𝗸 𝘁𝗼 𝗮 𝗛𝘂𝗺𝗮𝗻\n\n" +
		"For direct assistance, please reach out to us through the following channels. Our team is ready to help!\n\n" +
		"📞 **Phone:**\n(032) 123-4567\n\n" +
		"📧 **Email:**\nsales@goodkredit.com\n\n" +
		"🏢 **Office Hours:**\nMonday - Friday\n8:00 AM - 5:00 PM PHT"

	LinkAccountInstructions = "📝 𝗟𝗶𝗻𝗸 𝗮 𝗦𝘁𝘂𝗱𝗲𝗻𝘁 𝗔𝗰𝗰𝗼𝘂𝗻𝘁🧑‍🎓\n\n" +
		"Please contact your school administrator and share this unique code:\n\n" +
		"🔑 Your Account Code: %s\n\n" +
		"🧭 Once linked, you'll be able to view your grades, manage school fees, and more."

	AccountLinkedWelcome = "𝗔𝗰𝗰𝗼𝘂𝗻𝘁 𝗟𝗶𝗻𝗸𝗲𝗱!\n\n" +
		"Welcome, %s %s from %s!\n\n" +
		"You're all set. The School Assistant is now active for this profile."

	AccountPrimaryMessage = "👤 *%s %s* is currently your active profile.\n\n" +
		"Student ID: %s\n" +
		"School: %s\n\n" +
		"Would you like to continue with this profile or switch to another?"

	SingleProfileConfirmation = "I found this profile linked to your account:\n\n" +
		"👤 Name: %s %s\n" +
		"📚 Student ID: %s\n" +
		"🏫 School: %s\n\n" +
		"Would you like to continue with this profile?"

	WelcomeMessageAboard = "🎉 Welcome Aboard! 🎉\n\n" +
		"Your School Assistant is now ready to help you with:\n" +
		"• Viewing your grades and academic progress\n" +
		"• Managing school fees and payments\n" +
		"• Staying updated with school announcements\n" +
		"• Tracking your attendance\n\n" +
		"What would you like to do first?"

	MainMenuTemplate = "🏫 𝗠𝗮𝗶𝗻 𝗠𝗲𝗻𝘂\n\n" +
		"Student: %s\n" +
		"School: %s\n\n" +
		"Please choose an option:\n" +
		"[1] Check Grades\n" +
		"[2] School Fees\n" +
		"[3] School Bulletin\n" +
		"[4] View Attendance\n" +
		"[5] Manage Account\n" +
		"[6] Support"

	GradesTemplate = "📊 *Grades for %s*\n\n" +
		"*Subject* | *Grade* | *Semester*\n" +
		"-------------------------------\n" +
		"%s"

	GradeItemTemplate = "• %s | %s | %s\n"

	ProfileConfirmationMessage = "✅ Your active profile has been set to:\n\n👤 *%s %s*\nStudent ID: %s\nSchool: %s"
	AccountSelectionMessage    = "Select a profile to switch to:"
	AccountSwitchedMessage     = "✅ Successfully switched to %s %s's profile."
	NoPrimaryAccountMessage    = "No primary profile found. Please select a profile to continue:"

	// ProfileDetailsTemplate is the template for displaying user profile details
	ProfileDetailsTemplate = "📋Profile Details\n\n" +
		"Name: %s\n" +
		"Course: %s\n" +
		"Year Level: %s\n" +
		"Status: %s\n\n" +

		"Please choose an option:\n" +
		"[1] Subjects Enrolled\n" +
		"[2] Switch Profile\n"

	UserStatusUnregistered  = "UNREGISTERED"
	UserStatusRegistered    = "REGISTERED"
	UserStatusLinkedPrimary = "LINKED_PRIMARY"
	UserStatusDeactivated   = "DEACTIVATED"
)
