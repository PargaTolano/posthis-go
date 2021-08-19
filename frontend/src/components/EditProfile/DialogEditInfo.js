import React,{ Children, cloneElement } from 'react';

import { withStyles } from '@material-ui/core/styles';

import {
  Dialog,
  DialogTitle as MuiDialogTitle,
  DialogContent as MuiDialogContent,
  DialogActions as MuiDialogActions,
  IconButton,
  Typography
} from '@material-ui/core';

import {
  Close as CloseIcon,
  Settings as SettingsIcon
} from '@material-ui/icons';

import styles from '_styles/ProfileSettingsButton.module.css';

const dlgstyles = (theme) => ({
  root: {
    margin: 0,
    padding: theme.spacing(2),
  },
  closeButton: {
    position: 'absolute',
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
  
});

const DialogTitle = withStyles(dlgstyles)((props) => {
  const { children, classes, onClose, ...other } = props;
  return (
    <MuiDialogTitle disableTypography className={classes.root} {...other}>
      <Typography variant='h6'>{children}</Typography>
      {onClose ? (
        <IconButton aria-label='close' className={classes.closeButton} onClick={onClose}>
          <CloseIcon />
        </IconButton>
      ) : null}
    </MuiDialogTitle>
  );
});

const DialogContent = withStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
  },
}))(MuiDialogContent);

const DialogActions = withStyles((theme) => ({
  root: {
    margin: 0,
    padding: theme.spacing(1),
  },
}))(MuiDialogActions);

export const DialogEditInfo = ({children, color}) => {
  const [open, setOpen] = React.useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };

  const childrenWithProps = Children.map(children, (child, index) => {
    return cloneElement(child, {
      handleClose,
    });
  });

  return (
    <>
      <IconButton  
          aria-label='upload picture'
          component ='span' 
          variant   ='contained' 
          color     ={ color || 'secondary' } 
          onClick   ={handleClickOpen}
          className ={styles.btn}
      >
        <SettingsIcon />
      </IconButton>
      <Dialog style={{zIndex: 10001}} onClose={handleClose} aria-labelledby='customized-dialog-title' open={open}>
        <DialogContent dividers>
            {childrenWithProps}
        </DialogContent>
        
      </Dialog>
    </>
  );
};

export default DialogEditInfo