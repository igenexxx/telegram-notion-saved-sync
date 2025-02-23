To get a **Telegram App ID** and **App Hash**, you need to register your application with Telegram via their official developer portal at [my.telegram.org](https://my.telegram.org). This process is straightforward and requires a Telegram account tied to a phone number. Below are the step-by-step instructions to obtain these credentials:

---

### Steps to Get Telegram App ID and App Hash

1. **Visit Telegram’s Developer Portal**:
    - Open your web browser and go to [https://my.telegram.org](https://my.telegram.org).

2. **Log In with Your Phone Number**:
    - On the login page, enter the phone number associated with your Telegram account (e.g., `+1234567890`).
    - Click **Next**.
    - Telegram will send a verification code to your Telegram app (on your phone or desktop client). You won’t receive an SMS unless you’re logging in from a new device or session; instead, check your Telegram app for a message from "Telegram" with the code.

3. **Enter the Verification Code**:
    - Input the code you received in the Telegram app into the website’s prompt.
    - Click **Sign In**.

4. **Access the API Development Tools**:
    - Once logged in, you’ll see a dashboard with several options. Click on **API development tools** (or "Create application" if that’s the visible option).

5. **Fill Out the Application Form**:
    - You’ll be prompted to provide details about your application. Here’s what to enter:
        - **App title**: A name for your app (e.g., "MyTelegramSync").
        - **Short name**: A shorter version of the app name (e.g., "MTSync").
        - **URL**: Optional; you can leave it blank or enter a placeholder like "http://localhost".
        - **Platform**: Select the platform that best matches your use case (e.g., "Desktop" or "Other").
        - **Description**: Optional; briefly describe your app (e.g., "Syncs Saved Messages to Notion").
    - After filling out the form, click **Create application** (or "Save" depending on the interface).

6. **Retrieve App ID and App Hash**:
    - After submitting the form, you’ll see a page with your application details.
    - Look for:
        - **App api_id**: A numeric value (e.g., `123456`). This is your **App ID**.
        - **App api_hash**: A long alphanumeric string (e.g., `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6`). This is your **App Hash**.
    - Copy these values and store them securely. You’ll need them for your Go program’s configuration.

7. **Log Out (Optional)**:
    - Once you’ve saved your App ID and App Hash, you can log out of [my.telegram.org](https://my.telegram.org) by clicking your phone number in the top-right corner and selecting **Log out**.

---

### Example Config Usage
After obtaining your App ID and App Hash, add them to your `config.json` file like this:

```json
{
  "notion_token": "your_notion_integration_token",
  "telegram_app_id": 123456,
  "telegram_app_hash": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "telegram_phone": "+1234567890",
  "ai_key": "your_ai_api_key",
  "notion_database_id": "",
  "db_path": "local.db",
  "session_file": "telegram_session.json"
}
```

---

### Important Notes

- **Security**: Treat your App ID and App Hash as sensitive credentials. Do not share them publicly (e.g., in a public Git repository).
- **Phone Number**: The phone number must match an active Telegram account. Use international format (e.g., `+12025550123`).
- **Verification Code**: On your first login to [my.telegram.org](https://my.telegram.org) or when running your program with `gotd/td`, the code comes via Telegram messages, not SMS (unless SMS is explicitly requested by Telegram for security reasons).
- **Multiple Apps**: You can create multiple applications under the same account if needed, each with its own App ID and Hash.

---

### Troubleshooting
- **No Code Received**: Ensure you’re checking the Telegram app (not SMS) and that your phone number is correct.
- **Invalid Login**: Double-check your phone number format and ensure your Telegram account is active.
- **API Restrictions**: If you encounter issues, ensure your app isn’t violating Telegram’s terms (e.g., excessive requests).

Once you have your App ID and App Hash, you’re ready to use them in the Go program to authenticate with the Telegram 
MTProto API and access your "Saved Messages." Let me know if you run into any issues during this process!0