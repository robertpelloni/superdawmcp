import React from 'react';
import { render, fireEvent, act, waitFor } from '@testing-library/react-native';
import HomeScreen from '../src/screens/HomeScreen';
import TcpGateway from '../src/services/TcpGateway';

describe('HomeScreen MCP integration', () => {
  beforeEach(() => {
     TcpGateway.disconnect();
  });

  it('connects to gateway and toggles transport', async () => {
    const { getByText, findByText } = render(<HomeScreen />);
    const connectButton = getByText('Connect');

    act(() => {
      fireEvent.press(connectButton);
    });

    const connectedText = await findByText('Connected');
    expect(connectedText).toBeTruthy();

    const playButton = getByText('Play');

    act(() => {
      fireEvent.press(playButton);
    });

    await waitFor(() => {
        expect(getByText('Stop')).toBeTruthy();
    });
  });
});
