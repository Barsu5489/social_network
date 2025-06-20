<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { isAuthenticated, user, logout } from '../services/auth/authService';

  let authenticated = false;
  let currentUser = null;

  $: if ($isAuthenticated) {
    authenticated = true;
    currentUser = $user;
  }

  async function handleLogout() {
    await logout();
    goto('/login');
  }
</script>

<style>
  .navbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background-color: #f8f9fa;
    border-bottom: 1px solid #ccc;
  }

  .navbar a {
    margin: 0 1rem;
    text-decoration: none;
    color: #007bff;
    font-weight: bold;
  }

  .navbar a:hover {
    text-decoration: underline;
  }

  .user-info {
    display: flex;
    align-items: center;
  }

  .user-info span {
    margin-right: 0.5rem;
  }

  .logout-btn {
    background: none;
    border: none;
    color: #dc3545;
    cursor: pointer;
    font-weight: bold;
  }

  .logout-btn:hover {
    text-decoration: underline;
  }
</style>

  <nav class="navbar">
  <div>
    <a href="/">Home</a>
    {#if !authenticated}
      <a href="/login">Login</a>
      <a href="/register">Register</a>
    {/if}
  </div>
  {#if authenticated}
    <div class="user-info">
      <span>Welcome, {currentUser.first_name}</span>
      <button class="logout-btn" on:click={handleLogout}>
        Logout
      </button>
    </div>
  {/if}
</nav>
