
export interface User {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    nickname: string;
    date_of_birth: string;
    about_me: string;
    avatar_url: string;
    is_private: boolean;
    created_at: number;
}

export interface Post {
    id: string;
    user_id: string;
    group_id: string | null;
    content: string;
    privacy: 'public' | 'almost_private' | 'private';
    created_at: number;
    updated_at: number;
    deleted_at: number | null;
    likes_count: number;
    user_liked: boolean;
}

export interface Group {
    id: string;
    name: string;
    description: string;
    creator_id: string;
    is_private: boolean;
    created_at: number;
    updated_at: number;
}

export interface EventAttendee {
	id: string;
	event_id: string;
	user_id: string;
    user_name: string;
	status: 'going' | 'not_going' | 'maybe';
	joined_at: number;
}


export interface Event {
    id: string;
    group_id: string;
    title: string;
    description: string;
    location: string;
    start_time: number;
    end_time: number;
    created_by: string;
    created_at: number;
    updated_at: number;
    attendee_count: number;
    attendees: EventAttendee[];
}
