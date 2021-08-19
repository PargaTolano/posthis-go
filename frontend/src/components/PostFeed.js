import React, { useState, useEffect } from 'react';

import { NavBar, PostContainer }      from 'components/Feed';

import { CreatePostForm }             from 'components/Post';

import { PaginationElement }          from 'components/Common';

import { getFeed }                    from '_api';

import { loadingState }               from '_hooks';

import { toastService }               from '_services';

import styles from '_styles/PostFeed.module.css';

const limit = 4;

export const Feed = (props) => {

  const { history } = props;

  const [paginationActive, setPaginationActive] = useState(false);
  const [posts, setPosts] = useState([]);
  const [total, setTotal] = useState(0);
  const [last, setLast]   = useState(0);
  const [hasFetched, setHasFetched] = useState(false);
  const [refresh, setRefresh]       = useState(true);

  const toggleRefresh = ()=>setRefresh(x=>!x);

  const loadFeed = async ()=>{

    const {data:responseData, err} = await getFeed(0, Math.max(total, limit));
    if ( err !== null )
      return;

    const { data } = responseData;

    setHasFetched(true);

    if ( data === null) return;

    setLast( data.length );
    setTotal( x => x + data.length)
    setPosts(data);
    setPaginationActive(true);
  };

  useEffect(() => {
    loadFeed();
  },[]);

  useEffect(() => {
    (async()=>{

      if (!paginationActive || last < limit) return;

      const {data:responseData, err} = await getFeed(total, limit);
      loadingState.set(false);
      if ( err !== null ) {
        toastService.makeToast('Error on feed pagination call: ' + err, 'error');
        return;
      }
  
      const { data } = responseData;

      if ( data === null) return;
      
      setLast(data.length);
      setTotal( x => x + data.length   );
      setPosts( x => [...x, ...data] );
    })();
  },[refresh]);

  return (
    <div className={styles.root}>
      <NavBar history={history}/>
      <CreatePostForm afterUpdate={loadFeed}/>
      {
        (posts?.length !== 0) && <PostContainer posts={posts} history={history}/>
      }
      <PaginationElement
        name            = 'posts'
        hasFetched      = {hasFetched}
        total           = {total}
        last            = {last}
        limit           = {limit}
        onIntersection  = {toggleRefresh}
        rootMargin      = '400px'
      />
    </div>
  );
};

export default Feed;
