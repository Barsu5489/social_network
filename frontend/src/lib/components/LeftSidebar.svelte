<script>
  import { createPost, refreshPosts } from '$lib/services/postService'; // Import API functions
  export let user = null;
  
  const shortcuts = [
    { icon: '🧑‍🤝‍🧑', label: 'Friends', href: '/friends' },
    { icon: '👥', label: 'Groups', href: '/groups' },
  ];
  
  let showMore = false;
  let postText = '';
  let postVisibility = 'public';
  let errorMessage = '';
  let isSubmitting = false;
  
  function toggleShowMore() {
    showMore = !showMore;
  }
  
  async function submitPost() {
    if (!postText.trim()) return; // Require content
    
    isSubmitting = true;
    errorMessage = '';
    
    try {
      const postData = {
        content: postText,
        privacy: postVisibility,
        // Add groupID or allowedUserIDs if needed
        // groupID: someGroupID,
        // allowedUserIDs: someUserIDs,
      };
      
      const response = await createPost(postData);
      
      // Refresh posts to update feeds
      await refreshPosts();
      
      // Reset form
      postText = '';
    } catch (error) {
      // Parse error from backend if possible
      errorMessage = error.message.includes('401') ? 'Please log in to create a post.' :
                    error.message.includes('400') ? 'Invalid post content.' :
                    'Failed to create post. Please try again.';
      console.error('Post submission error:', error);
    } finally {
      isSubmitting = false;
    }
  }
</script>

<aside class="left-sidebar">
  <div class="sidebar-content">
    <!-- Post Creation Component -->
    <div class="post-composer">
      <div class="composer-header">
        <div class="user-avatar">
          {user?.name?.charAt(0) || 'U'}
        </div>
        <textarea 
          bind:value={postText}
          placeholder="What's on your mind?"
          class="post-input"
          rows="3"
          disabled={isSubmitting}
        ></textarea>
      </div>
      
      {#if errorMessage}
        <div class="error-message">{errorMessage}</div>
      {/if}
      
      <div class="composer-actions">
        <div class="left-actions">
          <select bind:value={postVisibility} class="visibility-select" disabled={isSubmitting}>
            <option value="public">🌐 Public</option>
            <option value="private">🔒 Private</option>
          </select>
        </div>
        
        <button 
          class="post-button" 
          on:click={submitPost}
          disabled={isSubmitting || !postText.trim()}
        >
          {isSubmitting ? 'Posting...' : 'Post'}
        </button>
      </div>
    </div>
    
    <div class="sidebar-divider"></div>
    
    <!-- User Profile -->
    <a href="/profile" class="sidebar-item user-item">
      <div class="item-icon profile-pic">
        {user?.name?.charAt(0) || 'U'}
      </div>
      <span class="item-label">{user?.name || 'User Name'}</span>
    </a>
    
    <!-- Main Navigation -->
    {#each shortcuts.slice(0, showMore ? shortcuts.length : 6) as shortcut}
      <a href={shortcut.href} class="sidebar-item">
        <div class="item-icon">
          {shortcut.icon}
        </div>
        <span class="item-label">{shortcut.label}</span>
      </a>
    {/each}
    
    <!-- Show More/Less Button -->
    <button class="sidebar-item show-more-btn" on:click={toggleShowMore}>
      <div class="item-icon">
        {showMore ? '🔼' : '🔽'}
      </div>
      <span class="item-label">{showMore ? 'See less' : 'See more'}</span>
    </button>
    
    <div class="sidebar-divider"></div>
    
    <!-- Your shortcuts section -->
    <div class="shortcuts-section">
      <h3 class="section-title">Your shortcuts</h3>
      <a href="/games/angry-birds" class="sidebar-item">
        <div class="item-icon game-icon">
          🐦
        </div>
        <span class="item-label">Angry Birds Friends</span>
      </a>
    </div>
  </div>
</aside>
<style>
  /* Post Composer Styles */
  .error-message {
    color: red;
    font-size: 0.9em;
    margin-top: 0.5em;
  }
  .post-composer {
    background: white;
    border: 1px solid #dadde1;
    border-radius: 8px;
    padding: 12px;
    margin-bottom: 16px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }
  
  .composer-header {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
  }
  
  .user-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: #1877f2;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 14px;
    flex-shrink: 0;
  }
  
  .post-input {
    flex: 1;
    border: none;
    outline: none;
    resize: none;
    font-size: 14px;
    font-family: inherit;
    background: #f0f2f5;
    border-radius: 20px;
    padding: 8px 12px;
    min-height: 36px;
    line-height: 1.4;
  }
  
  .post-input::placeholder {
    color: #65676b;
  }
  
  .composer-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 8px;
  }
  
  .left-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  
  .visibility-select {
    background: #f0f2f5;
    border: none;
    border-radius: 6px;
    padding: 4px 8px;
    font-size: 12px;
    cursor: pointer;
    color: #1c1e21;
  }
  
  .visibility-select:focus {
    outline: 2px solid #1877f2;
    outline-offset: 1px;
  }
  
  .post-button {
    background: #1877f2;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 6px 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .post-button:hover:not(:disabled) {
    background: #166fe5;
  }
  
  .post-button:disabled {
    background: #e4e6ea;
    color: #bcc0c4;
    cursor: not-allowed;
  }

  /* Original Sidebar Styles */
  .left-sidebar {
    position: fixed;
    left: 0;
    top: 56px;
    width: 280px;
    height: calc(100vh - 56px);
    background: white;
    overflow-y: auto;
    padding: 16px 0;
    border-right: 1px solid #dadde1;
  }
  
  .sidebar-content {
    padding: 0 8px;
  }
  
  .sidebar-item {
    display: flex;
    align-items: center;
    padding: 8px;
    border-radius: 8px;
    text-decoration: none;
    color: #1c1e21;
    transition: background-color 0.2s;
    margin-bottom: 4px;
    cursor: pointer;
    border: none;
    background: none;
    width: 100%;
    text-align: left;
    font-size: 15px;
  }
  
  .sidebar-item:hover {
    background: #f0f2f5;
  }
  
  .user-item .item-icon {
    background: #1877f2;
    color: white;
  }
  
  .item-icon {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 12px;
    font-size: 18px;
    flex-shrink: 0;
  }
  
  .profile-pic {
    font-weight: bold;
  }
  
  .game-icon {
    background: #f0f2f5;
    border-radius: 8px;
  }
  
  .item-label {
    font-weight: 500;
    color: #1c1e21;
  }
  
  .show-more-btn .item-label {
    color: #65676b;
  }
  
  .sidebar-divider {
    height: 1px;
    background: #dadde1;
    margin: 16px 8px;
  }
  
  .shortcuts-section {
    padding: 0 8px;
  }
  
  .section-title {
    font-size: 17px;
    font-weight: 600;
    color: #65676b;
    margin: 16px 0 8px 0;
    padding-left: 8px;
  }
  
  /* Scrollbar styling */
  .left-sidebar::-webkit-scrollbar {
    width: 8px;
  }
  
  .left-sidebar::-webkit-scrollbar-track {
    background: transparent;
  }
  
  .left-sidebar::-webkit-scrollbar-thumb {
    background: #bcc0c4;
    border-radius: 4px;
  }
  
  .left-sidebar::-webkit-scrollbar-thumb:hover {
    background: #8a8d91;
  }
</style>