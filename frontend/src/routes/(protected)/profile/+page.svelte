<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { user, isAuthenticated, logout } from '$lib/services/authService';

  let profileData = null;
  let loading = true;
  let error = "";

  const handleLogout = async () => {
    try {
      await logout();
      goto('/');
    } catch (err) {
      console.error('Logout failed:', err);
    }
  };

  const fetchProfile = async () => {
    try {
      const response = await fetch('/api/profile', {
        credentials: 'include' // Include cookies for session
      });
      
      if (response.ok) {
        const data = await response.json();
        profileData = data.data;
      } else {
        error = "Failed to load profile data";
      }
    } catch (err) {
      error = "Network error loading profile";
    } finally {
      loading = false;
    }
  };

  onMount(() => {
    // Check if user is authenticated
    const unsubscribe = isAuthenticated.subscribe(authenticated => {
      if (!authenticated) {
        goto('/login');
      } else {
        fetchProfile();
      }
    });

    return unsubscribe;
  });
</script>

<svelte:head>
  <title>Profile - ConnectSphere</title>
</svelte:head>

{#if $isAuthenticated}
  <div class="profile-container">
    <nav class="profile-nav">
      <h1 class="logo">ConnectSphere</h1>
      <button on:click={handleLogout} class="logout-btn">Logout</button>
    </nav>

    <main class="profile-main">
      {#if loading}
        <div class="loading">
          <div class="spinner"></div>
          <p>Loading your profile...</p>
        </div>
      {:else if error}
        <div class="error-state">
          <div class="error-icon">⚠️</div>
          <h2>Something went wrong</h2>
          <p>{error}</p>
          <button on:click={fetchProfile} class="retry-btn">Try Again</button>
        </div>
      {:else if profileData}
        <div class="profile-content">
          <div class="profile-header">
            <div class="avatar">
              <div class="avatar-placeholder">
                {profileData.user?.first_name?.charAt(0) || '?'}{profileData.user?.last_name?.charAt(0) || ''}
              </div>
            </div>
            <div class="profile-info">
              <h1>{profileData.user?.first_name || ''} {profileData.user?.last_name || ''}</h1>
              <p class="username">@{profileData.user?.username || 'username'}</p>
              <p class="email">{profileData.user?.email || ''}</p>
              {#if profileData.user?.about_me}
                <p class="bio">{profileData.user.about_me}</p>
              {/if}
            </div>
            <div class="profile-stats">
              <div class="stat">
                <span class="stat-number">{profileData.follower_count || 0}</span>
                <span class="stat-label">Followers</span>
              </div>
              <div class="stat">
                <span class="stat-number">{profileData.following_count || 0}</span>
                <span class="stat-label">Following</span>
              </div>
              <div class="stat">
                <span class="stat-number">{profileData.posts?.length || 0}</span>
                <span class="stat-label">Posts</span>
              </div>
            </div>
          </div>

          <div class="profile-sections">
            <div class="section">
              <h2>Recent Posts</h2>
              {#if profileData.posts && profileData.posts.length > 0}
                <div class="posts-grid">
                  {#each profileData.posts as post}
                    <div class="post-card">
                      <p>{post.content || 'Post content'}</p>
                      <small>Posted on {new Date(post.created_at).toLocaleDateString()}</small>
                    </div>
                  {/each}
                </div>
              {:else}
                <div class="empty-state">
                  <p>No posts yet. Start sharing your thoughts!</p>
                  <button class="create-post-btn">Create First Post</button>
                </div>
              {/if}
            </div>

            <div class="section">
              <h2>Connections</h2>
              <div class="connections">
                {#if profileData.followers && profileData.followers.length > 0}
                  <div class="connection-group">
                    <h3>Followers</h3>
                    <div class="users-list">
                      {#each profileData.followers.slice(0, 5) as follower}
                        <div class="user-item">
                          <div class="mini-avatar">{follower.first_name?.charAt(0) || '?'}</div>
                          <span>{follower.first_name} {follower.last_name}</span>
                        </div>
                      {/each}
                      {#if profileData.followers.length > 5}
                        <p class="more-count">+{profileData.followers.length - 5} more</p>
                      {/if}
                    </div>
                  </div>
                {/if}

                {#if profileData.following && profileData.following.length > 0}
                  <div class="connection-group">
                    <h3>Following</h3>
                    <div class="users-list">
                      {#each profileData.following.slice(0, 5) as following}
                        <div class="user-item">
                          <div class="mini-avatar">{following.first_name?.charAt(0) || '?'}</div>
                          <span>{following.first_name} {following.last_name}</span>
                        </div>
                      {/each}
                      {#if profileData.following.length > 5}
                        <p class="more-count">+{profileData.following.length - 5} more</p>
                      {/if}
                    </div>
                  </div>
                {/if}

                {#if (!profileData.followers || profileData.followers.length === 0) && (!profileData.following || profileData.following.length === 0)}
                  <div class="empty-state">
                    <p>Start connecting with people to build your network!</p>
                  </div>
                {/if}
              </div>
            </div>
          </div>
        </div>
      {:else}
        <div class="empty-state">
          <h2>Welcome to ConnectSphere!</h2>
          <p>Your profile is ready. Start connecting and sharing!</p>
        </div>
      {/if}
    </main>
  </div>
{:else}
  <div class="loading">
    <div class="spinner"></div>
    <p>Redirecting to login...</p>
  </div>
{/if}

<style>
  .profile-container {
    min-height: 100vh;
    background: #f8fafc;
  }

  .profile-nav {
    background: white;
    border-bottom: 1px solid #e5e7eb;
    padding: 1rem 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .logo {
    font-size: 1.75rem;
    font-weight: 700;
    background: linear-gradient(135deg, #6366f1, #ec4899);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin: 0;
  }

  .logout-btn {
    background: transparent;
    border: 2px solid #e5e7eb;
    color: #6b7280;
    padding: 0.5rem 1rem;
    border-radius: 0.5rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .logout-btn:hover {
    border-color: #dc2626;
    color: #dc2626;
  }

  .profile-main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }

  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 50vh;
    gap: 1rem;
  }

  .spinner {
    width: 2rem;
    height: 2rem;
    border: 3px solid #e5e7eb;
    border-top: 3px solid #6366f1;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .error-state {
    text-align: center;
    padding: 3rem;
  }

  .error-icon {
    font-size: 3rem;
    margin-bottom: 1rem;
  }

  .retry-btn {
    background: #6366f1;
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 0.5rem;
    font-weight: 600;
    cursor: pointer;
    margin-top: 1rem;
  }

  .profile-content {
    background: white;
    border-radius: 1rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
    overflow: hidden;
  }

  .profile-header {
    background: linear-gradient(135deg, #6366f1, #ec4899);
    color: white;
    padding: 2rem;
    display: grid;
    grid-template-columns: auto 1fr auto;
    gap: 2rem;
    align-items: center;
  }

  .avatar {
    position: relative;
  }

  .avatar-placeholder {
    width: 5rem;
    height: 5rem;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
    font-weight: 700;
    backdrop-filter: blur(10px);
  }

  .profile-info h1 {
    font-size: 2rem;
    font-weight: 700;
    margin: 0 0 0.5rem 0;
  }

  .username {
    opacity: 0.8;
    margin: 0 0 0.25rem 0;
  }

  .email {
    opacity: 0.7;
    font-size: 0.9rem;
    margin: 0 0 0.5rem 0;
  }

  .bio {
    opacity: 0.9;
    margin: 0;
  }

  .profile-stats {
    display: flex;
    gap: 2rem;
  }

  .stat {
    text-align: center;
  }

  .stat-number {
    display: block;
    font-size: 1.5rem;
    font-weight: 700;
  }

  .stat-label {
    font-size: 0.875rem;
    opacity: 0.8;
  }

  .profile-sections {
    padding: 2rem;
    display: grid;
    gap: 2rem;
  }

  .section h2 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1f2937;
    margin-bottom: 1rem;
  }

  .posts-grid {
    display: grid;
    gap: 1rem;
  }

  .post-card {
    background: #f9fafb;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    padding: 1rem;
  }

  .post-card p {
    margin: 0 0 0.5rem 0;
    color: #374151;
  }

  .post-card small {
    color: #6b7280;
  }

  .empty-state {
    text-align: center;
    padding: 2rem;
    color: #6b7280;
  }

  .create-post-btn {
    background: linear-gradient(135deg, #6366f1, #ec4899);
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 0.5rem;
    font-weight: 600;
    cursor: pointer;
    margin-top: 1rem;
  }

  .connections {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
  }

  .connection-group h3 {
    font-size: 1.125rem;
    font-weight: 600;
    color: #374151;
    margin-bottom: 1rem;
  }

  .users-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .user-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .mini-avatar {
    width: 2rem;
    height: 2rem;
    background: linear-gradient(135deg, #6366f1, #ec4899);
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 600;
  }

  .more-count {
    color: #6b7280;
    font-size: 0.875rem;
    margin: 0.5rem 0 0 0;
  }

  @media (max-width: 768px) {
    .profile-header {
      grid-template-columns: 1fr;
      text-align: center;
      gap: 1rem;
    }

    .profile-stats {
      justify-content: center;
    }

    .profile-main {
      padding: 1rem;
    }

    .profile-nav {
      padding: 1rem;
    }
  }
</style>
