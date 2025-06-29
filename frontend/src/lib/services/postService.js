import { writable } from 'svelte/store';
import { fetchWithSession } from '$lib/services/api';

// Store for posts data
export const posts = writable([]);
export const followingPosts = writable([]);

const API_URL = 'http://localhost:3000/api';

/**
 * Create a new post
 * @param {Object} postData - The post data
 * @param {string} postData.content - Post content
 * @param {string} postData.privacy - Post privacy setting
 * @param {string} [postData.groupID] - Optional group ID
 * @param {string[]} [postData.allowedUserIDs] - Optional array of allowed user IDs
 */
export async function createPost(postData) {
  try {
    const response = await fetchWithSession('/post', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        content: postData.content,
        privacy: postData.privacy,
        group_id: postData.groupID || null,
        allowed_user_ids: postData.allowedUserIDs || null
      })
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return response.json(); // Parse JSON response
  } catch (error) {
    console.error('Error creating post:', error);
    throw error;
  }
}

/**
 * Get posts from users you're following
 */
export async function getFollowingPosts() {
  try {
    const data = await fetchWithSession('/followPosts');
    followingPosts.set(data);
    return data;
  } catch (error) {
    console.error('Error fetching following posts:', error);
    followingPosts.set([]);
    throw error;
  }
}

/**
 * Get all posts (that the user has permission to see)
 */
export async function getAllPosts() {
  try {
    const data = await fetchWithSession('/posts');
    posts.set(data);
    return data;
  } catch (error) {
    console.error('Error fetching all posts:', error);
    posts.set([]);
    throw error;
  }
}

/**
 * Delete a post by ID
 * @param {string} postId - The ID of the post to delete
 */
export async function deletePost(postId) {
  try {
    const response = await fetchWithSession(`/delPost/${postId}`, {
      method: 'DELETE'
    });

    // Refresh posts after successful deletion
    await Promise.all([
      getAllPosts().catch(() => {}), // Don't fail if one of these fails
      getFollowingPosts().catch(() => {})
    ]);

    return response;
  } catch (error) {
    console.error('Error deleting post:', error);
    throw error;
  }
}

/**
 * Refresh both post feeds
 */
export async function refreshPosts() {
  try {
    const [allPostsData, followingPostsData] = await Promise.allSettled([
      getAllPosts(),
      getFollowingPosts()
    ]);

    return {
      allPosts: allPostsData.status === 'fulfilled' ? allPostsData.value : [],
      followingPosts: followingPostsData.status === 'fulfilled' ? followingPostsData.value : []
    };
  } catch (error) {
    console.error('Error refreshing posts:', error);
    throw error;
  }
}