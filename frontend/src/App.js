import React, { useState, useEffect }                         from 'react';
import { BrowserRouter as Router, Route, Switch }  from 'react-router-dom';

import { toast, ToastContainer } from 'react-toastify';

import { 
  Feed, 
  Login, 
  NotFound, 
  PostDetail, 
  ProfileDetail, 
  SearchResult 
} from 'components';

import { PrivateRoute, PublicRoute }              from 'components/Routing';
import { Loading }                                from 'components/Common';

import { routes }                                 from '_utils';
import { history }                                from '_helpers';
import { authenticationService, toastService }                  from '_services';

import 'react-toastify/dist/ReactToastify.min.css';
import 'react-toastify/dist/ReactToastify.minimal.css';

import '_styles/ToastStyles.css';

function App() {

  const [ user, setUser ] = useState(null);

  useEffect(()=>{

    let toastSubs = 
        toastService
            .toast$
            .subscribe(({content, type}) => 
              {
                toast[type]( content, {
                  position:         'top-left',
                  autoClose:        3000,
                  hideProgressBar:  false,
                  closeOnClick:     true,
                  pauseOnHover:     true,
                  draggable:        true,
                  progress:         undefined,
                })
              }
            );

    let authSubs = 
          authenticationService
              .currentUser
              .subscribe(x=> void setUser(x));


    return ()=>{
      toastSubs.unsubscribe();
      authSubs.unsubscribe();
    };

  },[]);
  
  const temp = { history, user };

  return (
    <div className= 'App'>
      <Loading/>
      <ToastContainer
        position="top-left"
        autoClose={3000}
        hideProgressBar={false}
        newestOnTop
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
      />
      <Router>
        <Switch>
          <PrivateRoute exact path={routes.feed}          component={Feed}            {...temp}   />
          <PrivateRoute exact path={routes.postDetail}    component={PostDetail}      {...temp}   />
          <PrivateRoute       path={routes.searchResult}  component={SearchResult}    {...temp}   />
          <PrivateRoute exact path={routes.profile}       component={ProfileDetail}   {...temp}   />
          <PublicRoute  exact path={routes.login}         component={Login}           {...temp}   />
          <Route        exact path={'*'}                  component={NotFound}        {...temp}   />
        </Switch>
      </Router>
      
    </div>
  );
}

export default App;
