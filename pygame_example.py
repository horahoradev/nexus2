import pygame
import datetime
import esper
import pytmx
from pytmx import TiledImageLayer
from pytmx import TiledObjectGroup
from pytmx import TiledTileLayer
from pytmx.util_pygame import load_pygame
import logging
logger = logging.getLogger(__name__)


FPS = 60
RESOLUTION = 720, 480

TILEBOARD_WIDTH = 0
TILEBOARD_HEIGHT = 0

class TiledRenderer(object):
    """
    Super simple way to render a tiled map
    """

    def __init__(self, filename):
        tm = load_pygame(filename)

        # self.size will be the pixel size of the map
        # this value is used later to render the entire map to a pygame surface
        self.pixel_size = tm.width * tm.tilewidth, tm.height * tm.tileheight
        global TILEBOARD_HEIGHT, TILEBOARD_WIDTH
        TILEBOARD_HEIGHT = tm.tileheight
        TILEBOARD_WIDTH = tm .tilewidth
        self.tmx_data = tm

    def render_map(self, surface, offsets):
        """ Render our map to a pygame surface

        Feel free to use this as a starting point for your pygame app.
        This method expects that the surface passed is the same pixel
        size as the map.

        Scrolling is a often requested feature, but pytmx is a map
        loader, not a renderer!  If you'd like to have a scrolling map
        renderer, please see my pyscroll project.
        """

        # fill the background color of our render surface
        if self.tmx_data.background_color:
            surface.fill(pygame.Color(self.tmx_data.background_color))

        # iterate over all the visible layers, then draw them
        for layer in self.tmx_data.visible_layers:
            # each layer can be handled differently by checking their type

            if isinstance(layer, TiledTileLayer):
                self.render_tile_layer(surface, layer, offsets)

            elif isinstance(layer, TiledObjectGroup):
                self.render_object_layer(surface, layer, offsets)

            elif isinstance(layer, TiledImageLayer):
                self.render_image_layer(surface, layer, offsets)

    def render_tile_layer(self, surface, layer, offsets):
        """ Render all TiledTiles in this layer
        """
        # deref these heavily used references for speed
        tw = self.tmx_data.tilewidth
        th = self.tmx_data.tileheight
        surface_blit = surface.blit

        # iterate over the tiles in the layer, and blit them
        for x, y, image in layer.tiles():
            surface_blit(image, (offsets[0] + x * tw, offsets[1] + y * th))

    def render_object_layer(self, surface, layer, offsets):
        """ Render all TiledObjects contained in this layer
        """
        # deref these heavily used references for speed
        draw_lines = pygame.draw.lines
        surface_blit = surface.blit

        # these colors are used to draw vector shapes,
        # like polygon and box shapes
        rect_color = (255, 0, 0)

        # iterate over all the objects in the layer
        # These may be Tiled shapes like circles or polygons, GID objects, or Tiled Objects
        for obj in layer:
            logger.info(obj)

            # objects with points are polygons or lines
            if obj.image:
                # some objects have an image; Tiled calls them "GID Objects"
                surface_blit(obj.image, (obj.x - offsets[0], obj.y - offsets[1]))

            else:
                # use `apply_transformations` to get the points after rotation
                draw_lines(surface, rect_color, obj.closed, obj.apply_transformations(), 3)

    def render_image_layer(self, surface, layer, offsets):
        if layer.image:
            surface.blit(layer.image, (-offsets[0], -offsets[1]))



##################################
#  Define some Components:
##################################
class Velocity:
    def __init__(self, x=0.0, y=0.0):
        self.x = x
        self.y = y


class Renderable:
    def __init__(self, image, posx, posy, depth=0):
        self.image = image
        self.depth = depth
        self.x = posx
        self.y = posy
        self.w = image.get_width()
        self.h = image.get_height()

class PlayerControlled:
    def __init__(self):
        return


filename = "data/sewers.tmx"

################################
#  Define some Processors:
################################
class TileboardRenderer(esper.Processor):
    def __init__(self, window):
        self.renderer = None
        self.running = False
        self.dirty = False
        self.exit_status = 0
        self.load_map(filename)
        self.window = window

    def load_map(self, filename):
        """ Create a renderer, load data, and print some debug info
        """
        self.renderer = TiledRenderer(filename)

        logger.info("Objects in map:")
        for obj in self.renderer.tmx_data.objects:
            logger.info(obj)
            for k, v in obj.properties.items():
                logger.info("%s\t%s", k, v)

        logger.info("GID (tile) properties:")
        for k, v in self.renderer.tmx_data.tile_properties.items():
            logger.info("%s\t%s", k, v)

        logger.info("Tile colliders:")
        for k, v in self.renderer.tmx_data.get_tile_colliders():
            logger.info("%s\t%s", k, list(v))

    def process(self):
        """ Draw our map to some surface (probably the display)
        """
        for ent, (vel, rend, pc) in self.world.get_components(Velocity, Renderable, PlayerControlled):
            # first we make a temporary surface that will accommodate the entire
            # size of the map.
            # because this demo does not implement scrolling, we render the
            # entire map each frame
            temp = pygame.Surface(self.renderer.pixel_size)

            # render the map onto the temporary surface
            self.renderer.render_map(temp, (-rend.x, -rend.y))

            # now resize the temporary surface to the size of the display
            # this will also 'blit' the temp surface to the display
            # pygame.transform.(temp, self.window.get_size(), self.window)
            self.window.blit(temp, (0, 0))


            # display a bit of use info on the display
            # f = pygame.font.Font(pygame.font.get_default_font(), 20)
            # i = f.render('press any key for next map or ESC to quit',
            #              1, (180, 180, 0))
            # self.window.blit(i, (0, 0))


class MovementProcessor(esper.Processor):
    def __init__(self, minx, maxx, miny, maxy):
        super().__init__()
        self.minx = minx
        self.maxx = maxx
        self.miny = miny
        self.maxy = maxy
        self.next_movement = datetime.datetime.now()

    def process(self):
        # This will iterate over every Entity that has BOTH of these components:
        if self.next_movement < datetime.datetime.now():
            for ent, (vel, rend, pc) in self.world.get_components(Velocity, Renderable, PlayerControlled):
                # Update the Renderable Component's position by it's Velocity:
                rend.x += vel.x * TILEBOARD_WIDTH
                rend.y += vel.y * TILEBOARD_WIDTH
                # An example of keeping the sprite inside screen boundaries. Basically,
                # adjust the position back inside screen boundaries if it tries to go outside:
                # rend.x = max(self.minx, rend.x)
                # rend.y = max(self.miny, rend.y)
                # rend.x = min(self.maxx - rend.w, rend.x)
                # rend.y = min(self.maxy - rend.h, rend.y)
                self.next_movement = self.next_movement + datetime.timedelta(0, 0.25)


class RenderProcessor(esper.Processor):
    def __init__(self, window, clear_color=(0, 0, 0)):
        super().__init__()
        self.window = window
        self.clear_color = clear_color

    def process(self):
        # Clear the window:
        # self.window.fill(self.clear_color)
        # This will iterate over every Entity that has this Component, and blit it:
        for ent, rend in self.world.get_component(Renderable):
            rend.image = pygame.transform.scale(rend.image, (TILEBOARD_WIDTH, TILEBOARD_HEIGHT))
            self.window.blit(rend.image, (RESOLUTION[0]/2, RESOLUTION[1]/2))
        # Flip the framebuffers
        pygame.display.flip()


################################
#  The main core of the program:
################################
def run():
    # Initialize Pygame stuff
    pygame.init()
    window = pygame.display.set_mode(RESOLUTION)
    pygame.display.set_caption("Esper Pygame example")
    clock = pygame.time.Clock()
    pygame.key.set_repeat(1, 1)

    # Initialize Esper world, and create a "player" Entity with a few Components.
    world = esper.World()
    player = world.create_entity()
    world.add_component(player, Velocity(x=0, y=0))
    world.add_component(player, Renderable(image=pygame.image.load("data/acid1.png"), posx=100, posy=100))
    world.add_component(player, PlayerControlled())
    # Another motionless Entity:
    enemy = world.create_entity()
    world.add_component(enemy, Renderable(image=pygame.image.load("data/acid1.png"), posx=400, posy=250))

    # Create some Processor instances, and asign them to be processed.
    render_processor = RenderProcessor(window=window)
    movement_processor = MovementProcessor(minx=0, maxx=RESOLUTION[0], miny=0, maxy=RESOLUTION[1])
    tileboard_render_processor = TileboardRenderer(window)

    world.add_processor(render_processor)
    world.add_processor(movement_processor)
    world.add_processor(tileboard_render_processor)

    running = True
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False
            elif event.type == pygame.KEYDOWN:
                if event.key == pygame.K_LEFT:
                    # Here is a way to directly access a specific Entity's
                    # Velocity Component's attribute (y) without making a
                    # temporary variable.
                    world.component_for_entity(player, Velocity).x = -1
                elif event.key == pygame.K_RIGHT:
                    # For clarity, here is an alternate way in which a
                    # temporary variable is created and modified. The previous
                    # way above is recommended instead.
                    player_velocity_component = world.component_for_entity(player, Velocity)
                    player_velocity_component.x = 1
                elif event.key == pygame.K_UP:
                    world.component_for_entity(player, Velocity).y = -1
                elif event.key == pygame.K_DOWN:
                    world.component_for_entity(player, Velocity).y = 1
                elif event.key == pygame.K_ESCAPE:
                    running = False
            elif event.type == pygame.KEYUP:
                if event.key in (pygame.K_LEFT, pygame.K_RIGHT):
                    world.component_for_entity(player, Velocity).x = 0
                if event.key in (pygame.K_UP, pygame.K_DOWN):
                    world.component_for_entity(player, Velocity).y = 0

        # A single call to world.process() will update all Processors:
        world.process()

        clock.tick(FPS)


if __name__ == "__main__":
    run()
    pygame.quit()
