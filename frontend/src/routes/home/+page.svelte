<script>
  import { onMount } from 'svelte';
  import { getAllPosts } from '../../lib/services/postService';
  import PostCard from '../../lib/components/PostCard.svelte';
  import Navbar from '../../lib/components/Navbar.svelte';
  import CreatePost from '../../lib/components/CreatePost.svelte';
  import LeftSidebar from '../../lib/components/LeftSidebar.svelte';
  import RightSidebar from '../../lib/components/RightSidebar.svelte';

  let posts = [];

  onMount(async () => {
    try {
      posts = await getAllPosts();
    } catch (error) {
      console.error('Error fetching posts:', error);
    }
  });

  function handlePostCreated(event) {
    const newPost = event.detail;
    posts = [newPost, ...posts];
  }
</script>

<Navbar />

<div class="page-wrapper">
  <div class="home-container">
    <div class="left-sidebar">
      <div class="sidebar-content">
        <LeftSidebar />
      </div>
    </div>
    <div class="main-content">
      <div class="create-post-wrapper">
        <CreatePost on:postCreated={handlePostCreated} />
      </div>
      <div class="posts-container">
        {#each posts as post (post.id)}
          <div class="post-wrapper">
            <PostCard {post} />
          </div>
        {/each}
      </div>
    </div>
    <div class="right-sidebar">
      <div class="sidebar-content">
        <RightSidebar />
      </div>
    </div>
  </div>
</div>

<style>
  .page-wrapper {
    min-height: 100vh;
    background: linear-gradient(135deg, var(--background-color) 0%, var(--background-color) 100%);
    background-attachment: fixed;
  }

  .home-container {
    display: grid;
    grid-template-columns: 250px minmax(500px, 800px) 300px;
    gap: 24px;
    padding: 80px 32px 32px 32px;
    max-width: 1400px;
    margin: 0 auto;
    position: relative;
  }

  .sidebar-content {
    background: var(--card-background);
    border-radius: 12px;
    padding: 16px;
    position: sticky;
    top: 90px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--border-color);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }

  .sidebar-content:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .main-content {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .create-post-wrapper {
    background: var(--card-background);
    border-radius: 12px;
    padding: 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--border-color);
    transition: transform 0.2s ease;
  }

  .create-post-wrapper:hover {
    transform: translateY(-2px);
  }

  .posts-container {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .post-wrapper {
    background: var(--card-background);
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--border-color);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }

  .post-wrapper:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  @media (max-width: 1200px) {
    .home-container {
      grid-template-columns: 200px minmax(400px, 1fr) 250px;
      padding: 80px 16px 16px 16px;
    }
  }

  @media (max-width: 992px) {
    .home-container {
      grid-template-columns: minmax(0, 1fr) 250px;
    }

    .left-sidebar {
      display: none;
    }
  }

  @media (max-width: 768px) {
    .home-container {
      grid-template-columns: 1fr;
      padding: 70px 12px 12px 12px;
    }

    .right-sidebar {
      display: none;
    }

    .create-post-wrapper,
    .post-wrapper {
      border-radius: 8px;
    }
  }
</style>
