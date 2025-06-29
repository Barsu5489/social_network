<script>
  import { createEventDispatcher } from "svelte";

  export let user = null;

  let searchQuery = "";
  const dispatch = createEventDispatcher();

  function handleSearch() {
    console.log("Searching for:", searchQuery);
  }

  function handleNotifications() {
    console.log("Opening notifications");
  }

  function handleMessages() {
    console.log("Opening messages");
  }

  function handleProfile() {
    console.log("Opening profile");
  }

  function handleCreatePost() {
    dispatch("createPost");
  }
</script>

<nav class="navbar">
  <div class="navbar-content">
    <!-- Left section -->
    <div class="navbar-left">
      <div class="logo">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="28"
          height="28"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="logo-icon"
        >
          <path
            d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.72"
          ></path>
          <path
            d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.72-1.72"
          ></path>
        </svg>
      </div>
      <div class="search-container">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="search-icon"
        >
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <input
          type="text"
          placeholder="Search"
          bind:value={searchQuery}
          on:keydown={(e) => e.key === "Enter" && handleSearch()}
          class="search-input"
        />
      </div>
    </div>

    <!-- Right section -->
    <div class="navbar-right">
      <button class="nav-tab active">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
          <polyline points="9 22 9 12 15 12 15 22"></polyline>
        </svg>
        <span>Homepage</span>
      </button>
      <button class="nav-tab">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
          <circle cx="9" cy="7" r="4"></circle>
          <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
        </svg>
        <span>Connections</span>
      </button>
      <button class="nav-tab" on:click={handleMessages}>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path
            d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
          ></path>
        </svg>
        <span>Messages</span>
      </button>
      <button
        class="nav-tab notification-btn"
        on:click={handleNotifications}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
          <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
        </svg>
        <span>Notifications</span>
        <span class="notification-badge">2</span>
      </button>
      <button class="nav-tab">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <circle cx="12" cy="12" r="3"></circle>
          <path
            d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
          ></path>
        </svg>
        <span>Tools</span>
      </button>
    </div>
  </div>
</nav>

<style>
  .navbar {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    background: var(--card-background);
    border-bottom: 1px solid var(--border-color);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    z-index: 1000;
    height: 64px;
  }

  .navbar-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    max-width: 1600px;
    margin: 0 auto;
    padding: 0 32px;
    height: 100%;
  }

  .navbar-left {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .logo {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 48px;
    height: 48px;
    background-color: var(--primary-purple);
    border-radius: 12px;
  }

  .logo-icon {
    color: white;
  }

  .search-container {
    position: relative;
  }

  .search-icon {
    position: absolute;
    left: 14px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-color-secondary);
  }

  .search-input {
    width: 240px;
    height: 40px;
    border: 1px solid var(--input-border);
    border-radius: 20px;
    background: var(--input-background);
    padding: 0 16px 0 44px;
    font-size: 15px;
    color: var(--text-color-primary);
    outline: none;
    transition: all 0.2s ease;
  }

  .search-input:focus {
    border-color: var(--primary-purple);
    box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.2);
  }

  .navbar-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .nav-tab {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: none;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    color: var(--text-color-secondary);
    font-weight: 500;
    font-size: 15px;
    transition: background-color 0.2s, color 0.2s;
  }

  .nav-tab:hover {
    background: var(--button-hover-background);
    color: var(--text-color-primary);
  }

  .nav-tab.active {
    color: var(--primary-purple);
    background: var(--light-purple-transparent, rgba(167, 139, 250, 0.1));
  }

  .notification-btn {
    position: relative;
  }

  .notification-badge {
    position: absolute;
    top: 6px;
    right: 6px;
    background: var(--pink-accent);
    color: white;
    border-radius: 50%;
    width: 8px;
    height: 8px;
    border: 2px solid var(--card-background);
  }
</style>