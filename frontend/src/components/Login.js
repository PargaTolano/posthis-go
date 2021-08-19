import React,{ useState, useEffect }from 'react';

import {
  Grid,
  TextField,
  FormControlLabel,
  Button,
  Link,
  CssBaseline,
  Avatar,
  Paper,
  Checkbox,
  Typography,
} from '@material-ui/core';

import { makeStyles }             from '@material-ui/core/styles';
import PersonPinIcon              from '@material-ui/icons/PersonPin';

import { DialogSignup, SignUp }   from 'components/Signup';

import { authenticationService, toastService }  from '_services';
import { routes }                 from '_utils';

import { loadingState }           from '_hooks';

import styles from '_styles/Login.module.css';

export const Login = (props) => {

  const { history} = props;

  const [ username, setUsername ] = useState('');
  const [ password, setPassword ] = useState('');

  const getOnChange = ( setState )=>e=>setState(e.target.value);

  const onSubmit = async e =>{
    e.preventDefault();

    loadingState.set(true);

    const data = await authenticationService
      .login( username, password )
    
    loadingState.set(false);

    if ( data === null)
      return;

    toastService.makeToast( data.message, 'success')
    history.push( routes.feed );
  };

  return (
    <Grid container component='main' className={styles.root}>
      <CssBaseline />
      <Grid item xs={false} sm={4} md={7} lg={8} className={styles.image} />
      <Grid item xs={12}    sm={8} md={5} lg={4} component={Paper} elevation={6} square>
        <div className={styles.paper}>
          <Avatar className={styles.avatar}>
            <PersonPinIcon />
          </Avatar>

          <Typography component='h1' variant='h5'>
            <strong>Log In</strong>
          </Typography>

          <form className={styles.form} noValidate onSubmit={onSubmit}>
            <TextField
              variant='outlined'
              margin='normal'
              required
              fullWidth
              id='email'
              label='Username'
              name='email'
              autoComplete='email'
              autoFocus
              onChange = {getOnChange(setUsername)}
            />

            <TextField
              variant='outlined'
              margin='normal'
              required
              fullWidth
              name='password'
              label='Password'
              type='password'
              id='password'
              autoComplete='current-password'
              onChange = {getOnChange(setPassword)}
            />

            <Button
              type='submit'
              fullWidth
              variant='contained'
              color='primary'
              className={styles.submit}
            >
              Log In
            </Button>
          </form>

          <div className={styles.modalContainer}>
            <DialogSignup>
              <SignUp/>
            </DialogSignup>
          </div>
          
        </div>
      </Grid>
    </Grid>
  );
};

export default Login;
