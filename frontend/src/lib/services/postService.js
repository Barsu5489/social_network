const API_URL = 'http://localhost:3000';

export async function getAllPosts() {
  const response = await fetch(`${API_URL}/posts`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include' // Send cookies with the request
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to fetch posts');
  }

  return response.json();
}

export async function getPost(postId) {
  const response = await fetch(`${API_URL}/posts/${postId}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include' // Send cookies with the request
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to fetch post');
  }

  return response.json();
}

export async function createPost(postData) {
  const headers = {};

  let body;
  if (postData instanceof FormData) {
    body = postData;
  } else {
    headers['Content-Type'] = 'application/json';
    body = JSON.stringify(postData);
  }

  const response = await fetch(`${API_URL}/post`, {
    method: 'POST',
    headers: headers,
    body: body,
    credentials: 'include' // Send cookies with the request
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to create post');
  }

  return response.json();
}

export async function likePost(postId) {
  const response = await fetch(`${API_URL}/posts/${postId}/like`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include' // Send cookies with the request
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to like post');
  }

  return response.json();
}

export async function unlikePost(postId) {
  const response = await fetch(`${API_URL}/posts/${postId}/like`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include' // Send cookies with the request
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to unlike post');
  }

  return response.json();
}

export async function getPostComments(postId) {
  const response = await fetch(`${API_URL}/comments/${postId}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include' // Send cookies with the request
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to fetch comments');
  }

  return response.json();
}

export async function addComment(postId, commentData) {
  const response = await fetch(`${API_URL}/comment/${postId}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(commentData),
    credentials: 'include' // Send cookies with the request
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to add comment');
  }

  return response.json();
}