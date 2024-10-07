import { GraphRoot } from '../graph-root';
import { GraphRenderer } from '../renderer';

export type SvgTestCallback<T> = (r: T) => void;

/**
 * Create a new graph root element for this spec
 * @param title Title of the spec
 * @param size Size of the wrapper div element. Size will be a square.
 * @param callback Test callback
 */
export function graphRootIt(
  title: string,
  size: number,
  callback: SvgTestCallback<GraphRoot>,
) {
  it(title, () => {
    const gr = generateRoot(title, size);
    callback(gr);
  });
}

export function rendererIt(
  title: string,
  size: number,
  callback: SvgTestCallback<GraphRenderer>,
) {
  it(title, () => {
    const gr = generateRenderer(title, size);
    callback(gr);
  });
}

function generateRoot(title: string, size = 300) {
  const gr = new GraphRoot();
  const svgWrap = generateWrapper(title, size);
  gr.attach(svgWrap);
  return gr;
}

function generateRenderer(title: string, size = 300) {
  const svgWrap = generateWrapper(title, size);
  return new GraphRenderer(svgWrap);
}

function generateWrapper(title: string, size = 300): HTMLDivElement {
  const wrapperDiv = document.createElement('div');
  const label = document.createElement('p');
  label.textContent = title;
  wrapperDiv.appendChild(label);
  wrapperDiv.style.border = '1px solid green';
  wrapperDiv.style.backgroundColor = '#DDD';
  const svgWrap = document.createElement('div');
  svgWrap.style.width = `80%`;
  svgWrap.style.height = `${size}px`;
  svgWrap.style.backgroundColor = 'white';
  wrapperDiv.appendChild(svgWrap);
  document.body.appendChild(wrapperDiv);
  return svgWrap;
}
