<script>
    let isLogin = true;
    let firstName = "";
    let lastName = "";
    let email = "";
    let password = "";
    let remember = false;
  
    import { register, login } from '$lib/services/auth/authService';

    const handleSubmit = async () => {
      if (isLogin) {
        try {
          const result = await login({ email, password });
          alert(`Login successful! ${JSON.stringify(result)}`);
        } catch (error) {
          alert(`Login failed: ${error.message}`);
        }
      } else {
        try {
          const result = await register({ firstName, lastName, email: email, password });
          alert(`Registration successful! ${JSON.stringify(result)}`);
        } catch (error) {
          alert(`Registration failed: ${error.message}`);
        }
      }
    };
  </script>
  
  <div class="auth-container">
    <div class="tab-toggle">
      <button on:click={() => isLogin = true} class:is-active={isLogin}>Login</button>
      <button on:click={() => isLogin = false} class:is-active={!isLogin}>Register</button>
    </div>
  
    <h2>{isLogin ? 'USER LOGIN' : 'REGISTER'}</h2>
  
    <form on:submit|preventDefault={handleSubmit}>
      <div class="input-group">
        <span class="icon">👤</span>
        <input type="text" placeholder="Enter first name" bind:value={firstName} required />
      </div>

      <div class="input-group">
        <span class="icon">👤</span>
        <input type="text" placeholder="Enter last name" bind:value={lastName} required />
      </div>
  
      {#if !isLogin}
        <div class="input-group">
          <span class="icon">📧</span>
          <input type="email" placeholder="Enter email" bind:value={email} required />
        </div>
      {/if}
  
      <div class="input-group">
        <span class="icon">🔒</span>
        <input type="password" placeholder="Enter password" bind:value={password} required minlength="6" />
      </div>
  
      {#if isLogin}
        <div class="options">
          <label><input type="checkbox" bind:checked={remember} /> Remember</label>
          <a href="#">Forgot password?</a>
        </div>
      {/if}
  
      <button type="submit" class="login-button">{isLogin ? 'LOGIN' : 'REGISTER'}</button>
    </form>
  </div>
