<script>
  import { register } from '../services/auth/authService';
  import { goto } from '$app/navigation';

  let email = '';
  let password = '';
  let firstName = '';
  let lastName = '';
  let error = '';

  async function handleRegister() {
    try {
      await register({
        email,
        password,
        firstName,
        lastName
      });
      goto('/login');
    } catch (err) {
      error = err.message;
    }
  }
</script>

<style>
  .register-container {
    max-width: 400px;
    margin: auto;
    padding: 2rem;
    border: 1px solid #ccc;
    border-radius: 8px;
  }

  .register-container h2 {
    margin-bottom: 1rem;
  }

  .register-container label {
    display: block;
    margin-bottom: 0.5rem;
  }

  .register-container input {
    width: 100%;
    padding: 0.5rem;
    margin-bottom: 1rem;
    border: 1px solid #ccc;
    border-radius: 4px;
  }

  .register-container button {
    width: 100%;
    padding: 0.75rem;
    background-color: #28a745;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .register-container button:hover {
    background-color: #218838;
  }

  .register-container .error {
    color: red;
    margin-bottom: 1rem;
  }
</style>

<div class="register-container">
  <h2>Register</h2>
  {#if error}
    <div class="error">{error}</div>
  {/if}
  <form on:submit|preventDefault={handleRegister}>
    <label for="firstName">First Name</label>
    <input
      type="text"
      id="firstName"
      bind:value={firstName}
      required
    />

    <label for="lastName">Last Name</label>
    <input
      type="text"
      id="lastName"
      bind:value={lastName}
      required
    />

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

    <button type="submit">Register</button>
  </form>
</div>
