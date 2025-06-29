import { redirect } from '@sveltejs/kit';
import { user as userStore } from '$lib/services/authService';

export async function load() {
  console.log('redirect in +layout.js:', redirect);
  try {
    const res = await getProfile();
    userStore.set(res.user); // hydrate Svelte store
    return { user: res.user };
  } catch (err) {
    throw redirect(302, '/login');
  }
}
