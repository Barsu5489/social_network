<script>
  import { goto } from '$app/navigation';
  
  let formData = {
    username: "",
    email: "",
    password: "",
    firstName: "",
    lastName: "",
    dateOfBirth: "",
    nickname: "",
    aboutMe: ""
  };
  let confirmPassword = "";
  let isLoading = false;
  let error = "";
  let step = 1;

  const handleNext = () => {
    if (!formData.username || !formData.email || !formData.password || !confirmPassword) {
      error = "Please fill in all required fields";
      return;
    }
    
    if (formData.password !== confirmPassword) {
      error = "Passwords do not match";
      return;
    }

    if (formData.password.length < 6) {
      error = "Password must be at least 6 characters long";
      return;
    }

    error = "";
    step = 2;
  };

  const handleBack = () => {
    step = 1;
    error = "";
  };

  const handleRegister = async () => {
    if (!formData.firstName || !formData.lastName || !formData.dateOfBirth) {
      error = "Please fill in all required fields";
      return;
    }

    isLoading = true;
    error = "";

    try {
      const response = await fetch('/api/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });
      
      if (response.ok) {
        const data = await response.json();
        // Registration successful, redirect to login
        goto('/login?registered=true');
      } else {
        const errorText = await response.text();
        error = errorText || 'Registration failed';
      }
    } catch (err) {
      error = 'Network error. Please try again.';
    } finally {
      isLoading = false;
    }
  };
</script>

<div class="register-container">
  <div class="register-card">
    <div class="register-header">
      <h1>Join ConnectSphere</h1>
      <p>Create your account to start connecting</p>
      <div class="progress-bar">
        <div class="progress" class:active={step >= 1}></div>
        <div class="progress" class:active={step >= 2}></div>
      </div>
    </div>

    {#if step === 1}
      <form on:submit|preventDefault={handleNext} class="register-form">
        <h3>Account Information</h3>
        
        {#if error}
          <div class="error-message">
            {error}
          </div>
        {/if}

        <div class="form-group">
          <label for="username">Username *</label>
          <input 
            id="username"
            type="text" 
            bind:value={formData.username}
            placeholder="Choose a username"
            required
          />
        </div>

        <div class="form-group">
          <label for="email">Email Address *</label>
          <input 
            id="email"
            type="email" 
            bind:value={formData.email}
            placeholder="Enter your email"
            required
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="password">Password *</label>
            <input 
              id="password"
              type="password" 
              bind:value={formData.password}
              placeholder="Create password"
              required
              minlength="6"
            />
          </div>

          <div class="form-group">
            <label for="confirm-password">Confirm Password *</label>
            <input 
              id="confirm-password"
              type="password" 
              bind:value={confirmPassword}
              placeholder="Confirm password"
              required
              minlength="6"
            />
          </div>
        </div>

        <button type="submit" class="next-btn">
          Next Step →
        </button>
      </form>
    {:else}
      <form on:submit|preventDefault={handleRegister} class="register-form">
        <h3>Personal Information</h3>
        
        {#if error}
          <div class="error-message">
            {error}
          </div>
        {/if}

        <div class="form-row">
          <div class="form-group">
            <label for="firstName">First Name *</label>
            <input 
              id="firstName"
              type="text" 
              bind:value={formData.firstName}
              placeholder="First name"
              required
            />
          </div>

          <div class="form-group">
            <label for="lastName">Last Name *</label>
            <input 
              id="lastName"
              type="text" 
              bind:value={formData.lastName}
              placeholder="Last name"
              required
            />
          </div>
        </div>

        <div class="form-group">
          <label for="dateOfBirth">Date of Birth *</label>
          <input 
            id="dateOfBirth"
            type="date" 
            bind:value={formData.dateOfBirth}
            required
          />
        </div>

        <div class="form-group">
          <label for="nickname">Nickname <span class="optional">(Optional)</span></label>
          <input 
            id="nickname"
            type="text" 
            bind:value={formData.nickname}
            placeholder="What should people call you?"
          />
        </div>

        <div class="form-group">
          <label for="aboutMe">About Me <span class="optional">(Optional)</span></label>
          <textarea 
            id="aboutMe"
            bind:value={formData.aboutMe}
            placeholder="Tell us about yourself..."
            rows="3"
          ></textarea>
        </div>

        <div class="button-row">
          <button type="button" on:click={handleBack} class="back-btn">
            ← Back
          </button>
          
          <button type="submit" class="register-btn" disabled={isLoading}>
            {#if isLoading}
              <span class="spinner"></span>
              Creating Account...
            {:else}
              Create Account
            {/if}
          </button>
        </div>
      </form>
    {/if}

    <div class="register-footer">
      <p>Already have an account? <a href="/login">Sign in here</a></p>
      <a href="/" class="back-link">← Back to home</a>
    </div>
  </div>
</div>

<style>
  .register-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 2rem;
  }

  .register-card {
    background: white;
    border-radius: 1rem;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
    padding: 3rem;
    width: 100%;
    max-width: 520px;
  }

  .register-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .register-header h1 {
    font-size: 2rem;
    font-weight: 700;
    color: #1f2937;
    margin-bottom: 0.5rem;
  }

  .register-header p {
    color: #6b7280;
    margin-bottom: 1.5rem;
  }

  .progress-bar {
    display: flex;
    gap: 0.5rem;
    justify-content: center;
  }

  .progress {
    width: 2rem;
    height: 0.25rem;
    background: #e5e7eb;
    border-radius: 0.125rem;
    transition: background-color 0.3s ease;
  }

  .progress.active {
    background: linear-gradient(135deg, #6366f1, #ec4899);
  }

  .register-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .register-form h3 {
    color: #374151;
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
  }

  .error-message {
    background: #fef2f2;
    border: 1px solid #fecaca;
    color: #dc2626;
    padding: 0.75rem;
    border-radius: 0.5rem;
    font-size: 0.875rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .form-group label {
    font-weight: 600;
    color: #374151;
    font-size: 0.875rem;
  }

  .optional {
    font-weight: 400;
    color: #9ca3af;
  }

  .form-group input,
  .form-group textarea {
    padding: 0.875rem;
    border: 2px solid #e5e7eb;
    border-radius: 0.5rem;
    font-size: 1rem;
    transition: all 0.2s ease;
    font-family: inherit;
  }

  .form-group input:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: #6366f1;
    box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
  }

  .next-btn,
  .register-btn {
    background: linear-gradient(135deg, #6366f1, #ec4899);
    color: white;
    border: none;
    padding: 1rem;
    border-radius: 0.5rem;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
  }

  .next-btn:hover,
  .register-btn:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
  }

  .register-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
    transform: none;
  }

  .button-row {
    display: flex;
    gap: 1rem;
  }

  .back-btn {
    background: transparent;
    color: #6b7280;
    border: 2px solid #e5e7eb;
    padding: 1rem;
    border-radius: 0.5rem;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    flex: 0 0 auto;
  }

  .back-btn:hover {
    background: #f9fafb;
    border-color: #d1d5db;
  }

  .register-btn {
    flex: 1;
  }

  .spinner {
    width: 1rem;
    height: 1rem;
    border: 2px solid transparent;
    border-top: 2px solid currentColor;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .register-footer {
    text-align: center;
    margin-top: 2rem;
    padding-top: 2rem;
    border-top: 1px solid #e5e7eb;
  }

  .register-footer p {
    color: #6b7280;
    margin-bottom: 1rem;
  }

  .register-footer a {
    color: #6366f1;
    text-decoration: none;
    font-weight: 600;
  }

  .register-footer a:hover {
    text-decoration: underline;
  }

  .back-link {
    font-size: 0.875rem;
    color: #9ca3af;
  }

  @media (max-width: 580px) {
    .register-container {
      padding: 1rem;
    }

    .register-card {
      padding: 2rem 1.5rem;
    }

    .form-row {
      grid-template-columns: 1fr;
    }

    .button-row {
      flex-direction: column;
    }
  }
</style>
