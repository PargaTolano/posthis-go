export const routes = {
    feed:           '/',
    postDetail:     '/post/:id',
    searchResult:   '/search/:query?',
    profile:        '/profile/:id',
    login:          '/login',
    getPost:        id=>`/post/${id}`,
    getSearch:      query=>`/search/${query}`,
    getProfile:     id=>`/profile/${id}`,
};