

How to obtain the `GOOGLE_API_KEY` and `GOOGLE_CSE_ID` from the Google Cloud Console to use with your `google-search-agent` application.

### Understanding the Keys

  * **`GOOGLE_API_KEY`**: This key authenticates your application with Google's services. It tells Google that you have a valid project and permission to use their APIs.
  * **`GOOGLE_CSE_ID`**: This ID specifies the "Programmable Search Engine" you want to use. This allows you to customize the search results, for example, by limiting them to specific websites.

Here is how to get both:

### Step 1: Create a Google Cloud Project

Before you can create any credentials, you need a Google Cloud Project to house them.

1.  **Go to the Google Cloud Console:** [https://console.cloud.google.com/](https://console.cloud.google.com/)
2.  **Create a New Project:** If you don't have one already, click the project dropdown in the top navigation bar and select "New Project".
3.  **Name Your Project:** Give it a descriptive name (e.g., "My Search Agent") and click "Create".

### Step 2: Enable the Custom Search API

For your API key to work with the Custom Search Engine, you must enable the corresponding API for your project.

1.  **Navigate to the API Library:** In the Google Cloud Console, use the navigation menu (the "hamburger" icon ☰) to go to **APIs & Services \> Library**.
2.  **Search for the API:** In the search bar, type "Custom Search API" and press Enter.
3.  **Enable It:** Click on the "Custom Search API" from the results and then click the **Enable** button.

### Step 3: Create Your GOOGLE\_API\_KEY

Now you will generate the actual API key.

1.  **Go to Credentials:** In the navigation menu, go to **APIs & Services \> Credentials**.
2.  **Create a New API Key:** Click the **+ Create Credentials** button at the top of the page and select **API key**.
3.  **Copy Your Key:** A new window will pop up with your generated API key. Click the copy icon to copy it. This is the value you will use for `your_actual_api_key_here`.
4.  **Secure Your Key (Recommended):** Click the **Restrict Key** button. Under **API restrictions**, select "Restrict key" and from the dropdown, choose the "Custom Search API". This ensures your key can only be used for this specific purpose. Click **Save**.

You now have your `GOOGLE_API_KEY`.

### Step 4: Create Your GOOGLE\_CSE\_ID (Programmable Search Engine ID)

This ID comes from a different Google service but is linked to your Cloud project.

1.  **Go to the Programmable Search Engine page:** [https://programmablesearchengine.google.com/](https://programmablesearchengine.google.com/)
2.  **Create a New Search Engine:** Click the **Add** or **Get Started** button.
3.  **Configure Your Search Engine:**
      * **Name:** Give your search engine a name.
      * **What to search:** This is the most important part.
          * To search the entire web, select the option to **"Search the entire web"**.
          * To limit the search to specific sites, enter the site addresses (e.g., `www.mysite.com`).
      * Click **Create**.
4.  **Get the Search Engine ID:**
      * After creation, you will be taken to a "Congratulations" page. Click the **Customize** button.
      * On the control panel page, under the **Basics** tab, you will find the **Search engine ID**.
      * Click the **Copy to clipboard** button to get your ID. This is the value for `your_custom_search_engine_id_here`.

### Step 5: Configure Your `.env` File

You are now ready to complete the setup for your `google-search-agent`.

1.  **Open or Create the `.env` file** in the root of your `google-search-agent` project.
2.  **Paste your credentials** into the file like this:

<!-- end list -->

```
GOOGLE_API_KEY=AIzaSyB...your...actual...key...
GOOGLE_CSE_ID=a1b2c3d4e5f6g7h8i
```

Replace the placeholder values with the actual key and ID you just copied. You can now proceed with the rest of the steps in the setup guide (e.g., `make build` and `make run`).