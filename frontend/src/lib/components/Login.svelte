<script>
  import { login } from '../services/auth/authService';
  import { isAuthenticated } from '../services/auth/authService';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let email = '';
  let password = '';
  let error = '';

  async function handleLogin() {
    try {
      await login({ email, password });
      goto('/');
    } catch (err) {
      error = err.message;
    }
  }

  onMount(() => {
    if ($isAuthenticated) {
      goto('/');
    }
  });
</script>

<style>
  .login-container {
    max-width: 400px;
    margin: auto;
    padding: 2rem;
    border: 1px solid #ccc;
    border-radius: 8px;
  }

  .login-container h2 {
    margin-bottom: 1rem;
  }

  .login-container label {
    display: block;
    margin-bottom: 0.5rem;
  }

  .login-container input {
    width: 100%;
    padding: 0.5rem;
    margin-bottom: 1rem;
    border: 1px solid #ccc;
    border-radius: 4px;
  }

  .login-container button {
    width: 100%;
    padding: 0.75rem;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .login-container button:hover {
    background-color: #0056b3;
  }

  .login-container .error {
    color: red;
    margin-bottom: 1rem;
  }
</style>

<div class="login-container">
  <h2>Login</h2>
  {#if error}
    <div class="error">{error}</div>
  {/if}
  <form on:submit|preventDefault={handleLogin}>
    <label for="email">Email</label>
    <input
      type="email"
      id="email"
      bind:value={email}
      required
    />

    <label for="password">Password</label>
    <input
      type="password"
      id="password"
      bind:value={password}
      required
    />

    <button type="submit">Login</button>
  </form>
</div>
