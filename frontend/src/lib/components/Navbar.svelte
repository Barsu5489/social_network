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
