
<script>
  import { createPost } from '../../lib/services/postService';
  import { user } from '../../lib/services/auth/authService';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let postContent = '';
  let postPrivacy = 'public';
  let selectedFile = null;
  let previewUrl = null;
  let error = null;

  function handleFileSelect(e) {
    const file = e.target.files[0];
    if (file) {
      if (file.size > 10 * 1024 * 1024) { 
        error = 'File size must be less than 10MB';
        return;
      }
      selectedFile = file;
      previewUrl = URL.createObjectURL(file);
      error = null;
    }
  }

  async function handleCreatePost() {
    if (!postContent.trim() && !selectedFile) {
      alert('Post must have content or an image.');
      return;
    }

    const formData = new FormData();
    formData.append('content', postContent);
    formData.append('privacy', postPrivacy);
    if (selectedFile) {
      formData.append('image', selectedFile);
    }

    try {
      const newPost = await createPost(formData);
      postContent = '';
      selectedFile = null;
      previewUrl = null;
      dispatch('postCreated', newPost);
    } catch (error) {
      console.error('Error creating post:', error);
      alert('Failed to create post. Please try again.');
    }
  }
</script>

<div class="create-post-card">
  <div class="input-section">
    <img src={$user?.avatar_url || "/placeholder-avatar.png"} alt="User Avatar" class="avatar" />
    <textarea
      class="post-input"
      placeholder="Share something..."
      bind:value={postContent}
    ></textarea>
  </div>
  {#if previewUrl}
    <div class="image-preview">
      <img src={previewUrl} alt="Selected preview" />
      <button on:click={() => { previewUrl = null; selectedFile = null; }}>&times;</button>
    </div>
  {/if}
  <div class="actions-section">
    <div class="action-buttons">
      <label for="file-upload" class="action-btn">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
        <span>Image</span>
      </label>
      <input id="file-upload" type="file" on:change={handleFileSelect} accept="image/png, image/jpeg, image/gif" style="display: none;" />
      
    </div>
    <div class="post-controls">
      <select bind:value={postPrivacy} class="privacy-select">
        <option value="public">Public</option>
        <option value="almost_private">Followers</option>
        <option value="private">Private</option>
      </select>
      <button on:click={handleCreatePost} class="post-button">Post</button>
    </div>
  </div>
  {#if error}
    <p class="error">{error}</p>
  {/if}
</div>

<style>
  .create-post-card {
    background-color: var(--card-background);
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    padding: 20px;
    border: 1px solid var(--border-color);
  }

  .input-section {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 16px;
    margin-bottom: 12px;
  }

  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: #ccc;
  }

  .post-input {
    flex-grow: 1;
    border: none;
    background: transparent;
    font-size: 1.1em;
    color: var(--text-color-primary);
    resize: none;
    height: 40px;
    padding-top: 8px;
  }

  .post-input:focus {
    outline: none;
  }

  .image-preview {
    position: relative;
    margin-bottom: 12px;
  }

  .image-preview img {
    max-width: 100%;
    max-height: 250px;
    border-radius: 8px;
  }

  .image-preview button {
    position: absolute;
    top: 8px;
    right: 8px;
    background: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    border-radius: 50%;
    width: 28px;
    height: 28px;
    cursor: pointer;
    font-size: 18px;
    line-height: 1;
  }

  .actions-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .action-buttons {
    display: flex;
    gap: 16px;
  }

  .action-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    color: var(--text-color-secondary);
    font-weight: 500;
    font-size: 0.95em;
    padding: 8px;
    border-radius: 6px;
    transition: background-color 0.2s ease;
  }

  .action-btn:hover {
    background-color: var(--button-hover-background);
  }

  .post-controls {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .privacy-select {
    padding: 8px 12px;
    border: 1px solid var(--input-border);
    border-radius: 6px;
    background-color: var(--button-hover-background);
    color: var(--text-color-primary);
    font-weight: 500;
  }

  .post-button {
    background-color: var(--primary-purple);
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s ease;
  }

  .post-button:hover {
    background-color: var(--deep-purple);
  }

  .error {
    color: var(--pink-accent);
    font-size: 0.9em;
    margin-top: 10px;
  }
</style>
