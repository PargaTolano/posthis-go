import React                        from 'react';
import { Route, Redirect }          from 'react-router-dom';
import { authenticationService }    from '_services';
import {  routes }                  from '_utils';

export const PublicRoute = ({ component: Component, ...rest }) => {
    return <Route {...rest} render={props => (
        (authenticationService.currentUserValue)
            ? <Redirect to={{ pathname: routes.feed, state: { from: props.location } }} />
            : <Component {...props} />
    )} />
};

export default PublicRoute;