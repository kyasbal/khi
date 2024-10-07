import { Observable, Subject, filter, map } from 'rxjs';
import {
  KHIWindowPacket,
  WindowConnectionProvider,
} from './window-connector.service';

const KHI_APPLICATION_TOKEN = 'kubernetes-history-inspector';

/**
 * Packet data wrapper for BroadcastChannel
 * To verify if the data is sent from KHI
 */
interface BroadcastChannelPacketWrap {
  packet: KHIWindowPacket<unknown>;
  applicationToken: string;
}

/**
 * WindowConnectionProvider using BroadcastChannel
 */
export class BroadcastChannelWindowConnectionProvider
  implements WindowConnectionProvider
{
  private readonly channel: BroadcastChannel;

  private readonly messageReceiver: Subject<BroadcastChannelPacketWrap> =
    new Subject();

  constructor(channelName = KHI_APPLICATION_TOKEN) {
    this.channel = new BroadcastChannel(channelName);
    this.channel.addEventListener('message', (message) => {
      const data = message.data;
      this.messageReceiver.next(data);
    });
  }

  send(data: KHIWindowPacket<unknown>): void {
    this.channel.postMessage({
      packet: data,
      applicationToken: KHI_APPLICATION_TOKEN,
    } as BroadcastChannelPacketWrap);
  }

  receive(): Observable<KHIWindowPacket<unknown>> {
    return this.messageReceiver.pipe(
      filter(
        (packetWrap) => packetWrap.applicationToken === KHI_APPLICATION_TOKEN,
      ),
      map((packetWrap) => packetWrap.packet),
    );
  }
}

/**
 * WindowConnectionProvider used for tests.
 * This connection provider allows connecting to the other in the same frame
 */
export class InMemoryWindowConnectionProvider
  implements WindowConnectionProvider
{
  private readonly messageReceiver: Subject<KHIWindowPacket<unknown>> =
    new Subject();

  send(data: KHIWindowPacket<unknown>): void {
    setTimeout(() => {
      this.messageReceiver.next(data);
    }, 10);
  }

  receive(): Observable<KHIWindowPacket<unknown>> {
    return this.messageReceiver;
  }
}
