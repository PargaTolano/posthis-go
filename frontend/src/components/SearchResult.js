import React                  from 'react';

import { NavBar }             from 'components/Feed';

import {
  SearchUserCard, 
  SearchPostCard
} from 'components/Search';

import { DialogFollowWrapper } from 'components/Follow';

import { useMakeSearch } from '_hooks';

import styles                  from '_styles/SearchResult.module.css';
import { NavbarWrapper } from './Common';

export const SearchResult = ( props ) => {

  const {auth, match, history } = props;

  const { query } = match.params;

  const [ready, response] = useMakeSearch(query || '');
  
  return (
    <div className={styles.root}>

      <DialogFollowWrapper history={history}/>
      <NavbarWrapper auth={auth} history={history}/>

      <section className={styles.resultSection}>
        <h3 className={styles.resultTitle}>Users</h3>
        <div className={styles.resultCards}>
          {
            ready && ( response.users?.map( user =><SearchUserCard key={user.id} user={user} auth={auth}/>) )
          }
        </div>
        <h4 className={styles.seeMore}> See more..</h4>
      </section>

      <section className={styles.resultSection}>
        <h3 className={styles.resultTitle}>PosThis!</h3>
        <div className={styles.resultCards}>
          { 
            ready && ( response.posts?.map( post =><SearchPostCard key={post.postID} post={post}/>) )
          }
        </div>
        <h4 className={styles.seeMore}> See more..</h4>
      </section>
    </div>
  );
};

export default SearchResult;