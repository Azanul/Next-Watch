import React from 'react';
import styled from 'styled-components';

interface GoogleButtonProps {
  onClick: () => void;
}

const GoogleButton: React.FC<GoogleButtonProps> = ({ onClick }) => {
  return (
    <StyledButton onClick={onClick}>
      <IconWrapper>
        <GoogleIcon src="https://upload.wikimedia.org/wikipedia/commons/c/c1/Google_%22G%22_logo.svg" alt="Google logo" />
      </IconWrapper>
      <ButtonText>Sign in with Google</ButtonText>
    </StyledButton>
  );
};

const StyledButton = styled.button`
  display: flex;
  align-items: center;
  background-color: #4285f4;
  color: white;
  border: none;
  border-radius: 2px;
  padding: 0;
  font-size: 14px;
  font-weight: 500;
  font-family: 'Roboto', sans-serif;
  cursor: pointer;
  box-shadow: 0 3px 4px 0 rgba(0,0,0,.25);
  transition: box-shadow .3s;
  margin: auto;

  &:hover {
    box-shadow: 0 0 6px #4285f4;
  }

  &:active {
    background: #1669F2;
  }
`;

const IconWrapper = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 40px;
  height: 40px;
  border-radius: 2px;
  background-color: white;
`;

const GoogleIcon = styled.img`
  width: 18px;
  height: 18px;
`;

const ButtonText = styled.span`
  padding: 5px 5px 5px 5px;
  font-weight: 500;
`;

export default GoogleButton;