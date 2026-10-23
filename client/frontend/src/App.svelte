<script lang="ts">
  // Lógica da app aqui.
  // import { ImpetoService } from "../bindings/github.com/meneko/impeto/client";
  //

  import { onMount } from "svelte";
  import {AppsService, App, ComponentAction, ComponentTrigger, UIElement} from "../bindings/github.com/meneko/impeto/client"

  let canvas: HTMLCanvasElement
  let canvasContext: CanvasRenderingContext2D;

  //canvas.addEventListener(onclick, onmousedown, onload)

  let apps: App[] = [];
  let selectedApp: number;


  interface ElementTrigger {
      x: number,
      y: number,
      width: number,
      height: number,
      trigger: ComponentTrigger,
      action: ComponentAction
  }

  let elementsTriggers: ElementTrigger[] = [];

  onMount(async () => {
      apps = await AppsService.GetAll()

      const elementCanvas = document.getElementsByTagName('canvas')[0];

      if (!elementCanvas) {
          throw new Error("Canvas doesn't exist on the window!");
      }

      canvas = elementCanvas;

      const context = canvas.getContext("2d");

      if (!context) {
          throw new Error("Canvas didn't gave 2d context!");
      }
      canvasContext = context;
  })

  async function imageParser(source: string): Promise<HTMLImageElement> {
      const response = await fetch(source);

      const imageUrl = URL.createObjectURL(await response.blob());
      const image = new Image();

      image

      return new Promise((resolve, reject) => {
          image.addEventListener("load", () => {
              URL.revokeObjectURL(imageUrl);
              resolve(image)
          })

          image.addEventListener("error", () => {
              URL.revokeObjectURL(imageUrl);
              reject(new Error("Erro na conversão da imagem!"));
          })

          image.src = imageUrl;
      })
  }

  function drawHud() {
      console.log(selectedApp)
      if (typeof selectedApp != "number") { return }

      const hudUIElements = apps.find(app => app.id == selectedApp);

      if (!hudUIElements) {
          drawNoHud();
          return;
      }

      hudUIElements.modal.forEach((uiElement) => {
          switch (uiElement.component.type) {
              case "Background":
                drawBackground(uiElement, canvasContext)
                break;
              //case "Button":
                drawButton(uiElement, canvasContext)
                break;
              default:
                console.log("TODO COMPONENT: " + uiElement.component)
                break;
          }
      })
  }

  /*
    No futuro criar uma função para evitar a repetição do cidgo no drawHudError (!...image) e o awaitImageParser

  */

  async function drawBackground(backgroundElement: UIElement, canvasContext: CanvasRenderingContext2D) {
      const {image} = backgroundElement;

      if (!image) {
          drawHudError()
          return;
      }

      let backgroundToDraw = await imageParser(image.source);

      canvasContext.drawImage(backgroundToDraw, 0, 0, 200, 200);
  }

  async function drawButton(buttonElement: UIElement, canvasContext: CanvasRenderingContext2D) {
      const {component , image, layout} = buttonElement;

      if (!image || !layout || !component) {
          drawHudError()
          return;
      }

      let buttonToDraw = await imageParser(image.source);

      let {width, height} = layout;
      let {x, y, z} = layout.position

      canvasContext.drawImage(buttonToDraw, x, y, width, height)

      const {trigger, action} = component;

      if (!trigger || !action) {
          console.warn("Botão desenhado, mas sem funcionalidade atribuida");
          return;
      }

      elementsTriggers.push(
          {x, y, width, height, trigger, action}
      )
  }

  function drawNoHud() {
      console.log("Sem Hud")
  }

  function drawHudError() {
      console.log("Erro na criação do Hud")
  }
</script>

<main class="container">
  <!-- -- <h1 class="title">Impeto</h1> -->
  <ol>
      {#each apps as app}
          <button style="all: unset" on:click={() => {
                    selectedApp = app.id;
                    drawHud();
                }}><li>
              <div class='pill {app.id == selectedApp ? "selected-pill" : ""}'>
                  <div class='app-icon'>
                    {@html app.icon}
                  </div>
                  <span class='app-name'>{app.name}</span>
              </div>
          </li></button>
      {/each}
  </ol>

  <div id='separator'></div>

  <div id='hud'>
      <canvas width="200" height="200" id='canvas-hud'>
      </canvas>
  </div>
</main>

<style>
  .container {
    flex: 1 0 auto;
    height: 100vh;
    width: 100wh;
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: stretch;
    text-align: left;
  }
  ol {
      display: flex;
      flex-direction: column;
      row-gap: 5px;
      list-style: none;
      margin: 0;
      padding: 1px;
      padding-top: 3px;
      padding-left: 3px;
      width: 150px;

      overflow-y: auto;
      box-sizing: border-box;
  }
  #separator {
      width: 1px;
      background-color: #48484A;
  }
  #hud {
      flex: 1;
      box-sizing: border-box;

      overflow: hidden;
      height: 100%;
      width: 75%;
  }

  .pill {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      width: 90%;
      border-radius: 25px;

      padding: 5px 0px 5px 10px;
  }
  .pill:hover {
      background-color: #1e1e1f;
  }
  .selected-pill {
      background-color: #141415;
  }
  .selected-pill:hover {
      background-color: #141415;
  }
  .app-icon {
      width: 24px;
      height: 24px;
  }
  .app-icon :global(svg) {
      width: 100%;
      height: 100%;
  }
  .app-name {
      color: rgba(255, 255, 255, 0.75);
      font-weight: 600;
  }

  ol::-webkit-scrollbar {
       display: none;
       width: 0;
       background: transparent;
   }
</style>
