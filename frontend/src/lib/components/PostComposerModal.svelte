<script>
  import { createPost } from '../services/postService';
  import { createEventDispatcher } from 'svelte';

  export let showModal = false;
  let postContent = '';
  const dispatch = createEventDispatcher();

  async function handleSubmit() {
    if (postContent.trim() === '') return;

    try {
      const newPost = await createPost({ content: postContent });
      dispatch('postCreated', newPost);
      postContent = ''; // Clear the input
      showModal = false; // Close the modal
    } catch (error) {
      console.error('Error creating post:', error);
      // Handle error, e.g., show an error message to the user
    }
  }

  function closeModal() {
    showModal = false;
  }
</script>

{#if showModal}
  <div class="modal-overlay" on:click|self={closeModal}>
    <div class="modal-content">
      <div class="modal-header">
        <h2>Create New Post</h2>
        <button class="close-button" on:click={closeModal}>&times;</button>
      </div>
      <div class="modal-body">
        <textarea
          placeholder="What's on your mind?"
          bind:value={postContent}
          rows="5"
        ></textarea>
      </div>
      <div class="modal-footer">
        <button class="cancel-button" on:click={closeModal}>Cancel</button>
        <button class="post-button" on:click={handleSubmit}>Post</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal-content {
    background-color: var(--background-white);
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    width: 90%;
    max-width: 500px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px 20px;
    border-bottom: 1px solid var(--border-color);
  }

  .modal-header h2 {
    margin: 0;
    font-size: 1.2em;
    color: var(--text-primary);
  }

  .close-button {
    background: none;
    border: none;
    font-size: 1.5em;
    color: var(--text-gray);
    cursor: pointer;
  }

  .close-button:hover {
    color: var(--text-primary);
  }

  .modal-body {
    padding: 20px;
  }

  .modal-body textarea {
    width: 100%;
    padding: 10px;
    border: 1px solid var(--input-border);
    border-radius: 5px;
    font-size: 1em;
    resize: vertical;
    min-height: 100px;
    font-family: inherit;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    padding: 15px 20px;
    border-top: 1px solid var(--border-color);
    gap: 10px;
  }

  .cancel-button {
    background-color: #e0e0e0;
    color: var(--text-primary);
    border: none;
    padding: 8px 15px;
    border-radius: 5px;
    cursor: pointer;
  }

  .cancel-button:hover {
    background-color: #d0d0d0;
  }

  .post-button {
    background-color: var(--primary-purple);
    color: white;
    border: none;
    padding: 8px 15px;
    border-radius: 5px;
    cursor: pointer;
  }

  .post-button:hover {
    background-color: var(--deep-purple);
  }
</style>
